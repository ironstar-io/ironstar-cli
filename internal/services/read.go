package services

import (
	"io/ioutil"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	yaml "gopkg.in/yaml.v2"
)

func ReadInProjectConfig(projectRoot string) (types.ProjectConfig, error) {
	projConf := types.ProjectConfig{}
	py := filepath.Join(projectRoot, ".ironstar", "config.yml")

	err := SafeTouchConfigYAML(py)
	if err != nil {
		return projConf, err
	}

	pBytes, err := ioutil.ReadFile(py)
	if err != nil {
		return projConf, err
	}

	err = yaml.Unmarshal(pBytes, &projConf)
	if err != nil {
		return projConf, err
	}

	return projConf, nil
}

func ReadInCredentials() ([]types.Credentials, error) {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	err := SafeTouchConfigYAML(cp)
	if err != nil {
		return nil, err
	}

	cBytes, err := ioutil.ReadFile(cp)
	if err != nil {
		return nil, err
	}

	creds := []types.Credentials{}
	err = yaml.Unmarshal(cBytes, &creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
