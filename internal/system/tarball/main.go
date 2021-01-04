package tarball

import (
	"archive/tar"
	"compress/gzip"

	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewTarGZ walks path to create tar file tarName
func NewTarGZ(tarName string, path string, excludeFiles []string) (err error) {
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
		for _, excl := range excludeFiles {
			match, err := filepath.Match(excl, relFilePath)
			if err != nil {
				return err
			}

			umatch, err := filepath.Match(excl, unixFilePath)
			if err != nil {
				return err
			}

			if (match || umatch) && finfo.IsDir() {
				return filepath.SkipDir
			}

			if match || umatch {
				return nil
			}
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
