package services

import (
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	yaml "gopkg.in/yaml.v2"
)

func ReadInProjectConfig(projectRoot string) (types.ProjectConfig, error) {
	projConf := types.ProjectConfig{}
	py := filepath.Join(projectRoot, ".ironstar", "config.yml")

	err := SafeTouchConfigYAML(py)
	if err != nil {
		return projConf, err
	}

	pBytes, err := os.ReadFile(py)
	if err != nil {
		return projConf, err
	}

	err = yaml.Unmarshal(pBytes, &projConf)
	if err != nil {
		return projConf, err
	}

	return projConf, nil
}

func ReadInCredentials() (types.Credentials, error) {
	empty := types.Credentials{}
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	err := SafeTouchCredentialsYAML()
	if err != nil {
		return empty, err
	}

	cBytes, err := os.ReadFile(cp)
	if err != nil {
		return empty, err
	}

	creds := types.Credentials{}
	err = yaml.Unmarshal(cBytes, &creds)
	if err != nil {
		return empty, err
	}

	return creds, nil
}
