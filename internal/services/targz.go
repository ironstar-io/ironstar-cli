package services

import (
	"fmt"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/console"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/tarball"

	uuid "github.com/satori/go.uuid"
)

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

	fsplit := strings.Split(flg.Exclude, ",")
	excl := append(proj.Package.Exclude, fsplit...)

	fs.Mkdir("/tmp/ironstar")

	fmt.Println()
	wo := console.SpinStart("Creating a tarball containing your project files")

	tarpath := "/tmp/ironstar/" + uuid.NewV4().String() + ".tar.gz"
	err = tarball.NewTarGZ(tarpath, pr, excl)
	if err != nil {
		console.SpinPersist(wo, "⛔", "There was an error while creating a tarball for this project\n")
		return "", err
	}

	console.SpinPersist(wo, "🗜️", " A tarball containing your project files has been completed\n")

	return tarpath, nil
}
