package services

import (
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/tarball"

	uuid "github.com/satori/go.uuid"
)

// CreateProjectTar - Create a project tarball in tmp
func CreateProjectTar(exclFlag string) (string, error) {
	pr := fs.ProjectRoot()
	proj, err := ReadInProjectConfig(pr)
	if err != nil {
		return "", err
	}

	fsplit := strings.Split(exclFlag, ",")
	excl := append(proj.Package.Exclude, fsplit...)

	tarpath := "/tmp/" + uuid.NewV4().String() + ".tar.gz"
	err = tarball.NewTarGZ(tarpath, pr, excl)
	if err != nil {
		return "", err
	}

	return tarpath, nil
}
