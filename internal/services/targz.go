package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/system/console"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/tarball"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	uuid "github.com/satori/go.uuid"
)

// packageSizeWarnBytes is the built-tarball size above which we caution the
// user: large packages risk exceeding the API's upload window on slow links.
const packageSizeWarnBytes = 250 * 1024 * 1024

// CreateProjectTar - Create a project tarball in tmp
func CreateProjectTar(flg flags.Accumulator) (string, error) {
	if flg.CustomPackage != "" {
		return flg.CustomPackage, nil
	}

	pr := fs.ProjectRoot()
	proj, err := ReadInProjectConfig(pr)
	if err != nil {
		return "", err
	}

	ex, src, err := resolvePackageExcluder(pr, proj, flg)
	if err != nil {
		return "", err
	}

	fs.Mkdir("/tmp/ironstar")

	fmt.Println()
	wo := console.SpinStart("Creating a tarball containing your project files")

	tarpath := "/tmp/ironstar/" + uuid.NewV4().String() + ".tar.gz"
	err = tarball.NewTarGZ(tarpath, pr, ex)
	if err != nil {
		console.SpinPersist(wo, "⛔", "There was an error while creating a tarball for this project\n")
		return "", err
	}

	var size int64
	if fi, serr := os.Stat(tarpath); serr == nil {
		size = fi.Size()
	}

	console.SpinPersist(wo, "🗜️", fmt.Sprintf(" A tarball containing your project files has been completed (%s)\n", humanize.IBytes(uint64(size))))

	if size >= packageSizeWarnBytes {
		color.Yellow("⚠  This package is large (%s). Uploads may exceed the server's upload window on slower connections.\n   Review your excludes for database dumps, media, or build artifacts, or run `iron package --dry-run` to see what's included.", humanize.IBytes(uint64(size)))
	}

	printUploadNotice(src)

	return tarpath, nil
}

// WritePackageIndex lists the files that `iron package` would include (without
// uploading) and writes them to /tmp/ironstar-package-index-<datestamp>.txt.
// Sizes are uncompressed; the uploaded .tar.gz is smaller.
func WritePackageIndex(flg flags.Accumulator) (indexPath string, total int64, count int, err error) {
	var (
		entries      []tarball.IndexEntry
		source       string
		excludeLabel string
	)

	if flg.CustomPackage != "" {
		source = flg.CustomPackage
		excludeLabel = "n/a (custom package)"
		entries, total, err = tarball.IndexArchive(flg.CustomPackage)
	} else {
		source = fs.ProjectRoot()
		proj, rerr := ReadInProjectConfig(source)
		if rerr != nil {
			return "", 0, 0, rerr
		}
		ex, src, rerr := resolvePackageExcluder(source, proj, flg)
		if rerr != nil {
			return "", 0, 0, rerr
		}
		excludeLabel = src.label
		if excludeLabel == "" {
			excludeLabel = "none"
		}
		entries, total, err = tarball.Index(source, ex)
	}
	if err != nil {
		return "", 0, 0, err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Ironstar package dry-run index\n")
	fmt.Fprintf(&b, "# Generated: %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Fprintf(&b, "# Source: %s\n", source)
	fmt.Fprintf(&b, "# Excludes: %s\n", excludeLabel)
	fmt.Fprintf(&b, "# Files: %d   Total (uncompressed): %s\n", len(entries), humanize.IBytes(uint64(total)))
	fmt.Fprintf(&b, "# Note: sizes are uncompressed; the uploaded .tar.gz will be smaller.\n\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "%s\t%s\n", humanize.IBytes(uint64(e.Size)), e.Path)
	}

	indexPath = "/tmp/ironstar-package-index-" + time.Now().UTC().Format("20060102-150405") + ".txt"
	if err := fs.TouchByteArray(indexPath, []byte(b.String()), 0644); err != nil {
		return "", 0, 0, err
	}

	return indexPath, total, len(entries), nil
}
