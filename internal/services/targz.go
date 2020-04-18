package services

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/tarball"

	uuid "github.com/satori/go.uuid"
)

// CreateProjectTar - Create a project tarball in tmp
func CreateProjectTar() (string, error) {
	pr := fs.ProjectRoot()
	proj, err := ReadInProjectConfig(pr)
	if err != nil {
		return "", err
	}

	tarpath := "/tmp/" + uuid.NewV4().String() + ".tar.gz"
	err = tarball.NewTarGZ(tarpath, pr, proj.Package.Exclude)
	if err != nil {
		return "", err
	}

	return tarpath, nil
}
