package tarball

import (
	"archive/tar"
	"compress/gzip"

	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Excluder decides whether a walked path is omitted from a package. relPath is
// the slash-separated path relative to the package root. It is the single point
// both NewTarGZ and Index consult, so the packaged set and the dry-run index
// can never disagree.
type Excluder interface {
	Excludes(relPath string, isDir bool) (bool, error)
}

// PatternExcluder is the legacy config.yml `package.exclude` / `--exclude`
// matcher. A pattern containing a slash is matched against the whole relative
// path (root-anchored); a slash-less pattern also matches the basename at any
// depth, so `.git` or `node_modules` are excluded wherever they appear.
type PatternExcluder struct {
	patterns []string
}

func NewPatternExcluder(patterns []string) PatternExcluder {
	cleaned := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if p = strings.TrimSpace(p); p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return PatternExcluder{patterns: cleaned}
}

func (e PatternExcluder) Excludes(relPath string, isDir bool) (bool, error) {
	base := path.Base(relPath)
	for _, excl := range e.patterns {
		match, err := path.Match(excl, relPath)
		if err != nil {
			return false, fmt.Errorf("invalid exclude pattern %q: %w", excl, err)
		}
		if match {
			return true, nil
		}

		if !strings.Contains(excl, "/") {
			bmatch, err := path.Match(excl, base)
			if err != nil {
				return false, fmt.Errorf("invalid exclude pattern %q: %w", excl, err)
			}
			if bmatch {
				return true, nil
			}
		}
	}

	return false, nil
}

// NewTarGZ walks path to create tar file tarName, omitting anything ex excludes.
func NewTarGZ(tarName string, path string, ex Excluder) (err error) {
	tarFile, err := os.Create(tarName)
	if err != nil {
		return err
	}
	defer func() {
		err = tarFile.Close()
	}()

	absTar, err := filepath.Abs(tarName)
	if err != nil {
		return err
	}

	// enable compression if file ends in .gz
	tw := tar.NewWriter(tarFile)
	if strings.HasSuffix(tarName, ".gz") || strings.HasSuffix(tarName, ".gzip") {
		gz := gzip.NewWriter(tarFile)
		defer gz.Close()
		tw = tar.NewWriter(gz)
	}
	defer tw.Close()

	// validate path
	path = filepath.Clean(path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if absPath == absTar {
		return fmt.Errorf("tar file %s cannot be the source\n", tarName)
	}
	if absPath == filepath.Dir(absTar) {
		return fmt.Errorf("tar file %s cannot be in source %s\n", tarName, absPath)
	}

	walker := func(file string, finfo os.FileInfo, err error) error {
		// Return immediately if there's an error
		if err != nil {
			return err
		}

		var link string
		if finfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(file); err != nil {
				return err
			}
		}

		relFilePath, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}
		unixFilePath := strings.ReplaceAll(relFilePath, "\\", "/")

		// Don't include any files that are explicitly excluded by the user
		skip, err := ex.Excludes(unixFilePath, finfo.IsDir())
		if err != nil {
			return err
		}
		if skip {
			if finfo.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// fill in header info using func FileInfoHeader
		hdr, err := tar.FileInfoHeader(finfo, link)
		if err != nil {
			return err
		}

		// ensure header has relative file path (unix style)
		hdr.Name = unixFilePath

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// if path is a dir, dont continue
		if finfo.Mode().IsDir() {
			return nil
		}

		if !finfo.Mode().IsRegular() {
			return nil
		}

		// add file to tar
		srcFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(tw, srcFile)
		if err != nil {
			return err
		}

		return nil
	}

	// build tar
	if err := filepath.Walk(path, walker); err != nil {
		fmt.Printf("failed to add %s to tar: %s\n", path, err)
	}

	return nil
}

// IndexEntry is a regular file that would be included in a package, with its
// uncompressed size.
type IndexEntry struct {
	Path string
	Size int64
}

// Index walks path applying the same include/exclude rules as NewTarGZ and
// returns the regular files that would be packaged (in walk order) plus their
// total uncompressed size. It writes nothing.
func Index(path string, ex Excluder) (entries []IndexEntry, total int64, err error) {
	path = filepath.Clean(path)

	walker := func(file string, finfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relFilePath, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}
		unixFilePath := strings.ReplaceAll(relFilePath, "\\", "/")

		skip, err := ex.Excludes(unixFilePath, finfo.IsDir())
		if err != nil {
			return err
		}
		if skip {
			if finfo.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !finfo.Mode().IsRegular() {
			return nil
		}

		entries = append(entries, IndexEntry{Path: unixFilePath, Size: finfo.Size()})
		total += finfo.Size()

		return nil
	}

	if err := filepath.Walk(path, walker); err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// IndexArchive lists the regular-file entries of an existing .tar.gz, used when
// a pre-built custom package is supplied. Sizes are the uncompressed entry sizes.
func IndexArchive(tarGzPath string) (entries []IndexEntry, total int64, err error) {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, 0, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		entries = append(entries, IndexEntry{Path: hdr.Name, Size: hdr.Size})
		total += hdr.Size
	}

	return entries, total, nil
}
