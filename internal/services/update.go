package services

import (
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	yaml "gopkg.in/yaml.v2"
)

func UpdateCredentialsFile(newCreds types.Keylink) error {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	// Pull out matching login if it exists in the map
	var splicedKeychain []types.Keylink
	for _, cred := range credSet.Keychain {
		if cred.Login != newCreds.Login {
			splicedKeychain = append(splicedKeychain, cred)
		}
	}

	// Replace/Add the new credentials to the struct slice
	newKeychain := append(splicedKeychain, newCreds)

	newMarhsalled, err := yaml.Marshal(&types.Credentials{
		Active:   credSet.Active,
		Keychain: newKeychain,
	})
	if err != nil {
		return err
	}

	fs.Replace(cp, newMarhsalled)

	return nil
}

func RemoveCredentials(oldCreds types.Keylink) error {
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	// Pull out matching login if it exists in the map
	var splicedKeychain []types.Keylink
	for _, cred := range credSet.Keychain {
		if cred.Login != oldCreds.Login {
			splicedKeychain = append(splicedKeychain, cred)
		}
	}

	active := credSet.Active
	// Unset active login if removed set was active
	if credSet.Active == oldCreds.Login {
		active = ""
	}

	// Remarshal with the affected login credentials spliced out
	newMarhsalled, err := yaml.Marshal(&types.Credentials{
		Active:   active,
		Keychain: splicedKeychain,
	})
	if err != nil {
		return err
	}

	fs.Replace(cp, newMarhsalled)

	return nil
}

func UpdateActiveCredentials(login string) (types.Keylink, error) {
	empty := types.Keylink{}
	cp := filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml")

	credSet, err := ReadInCredentials()
	if err != nil {
		return empty, err
	}

	var credMatch types.Keylink
	for _, cred := range credSet.Keychain {
		if cred.Login == login {
			credMatch = cred
		}
	}
	if credMatch == (types.Keylink{}) {
		return empty, errs.NoCredentialMatch
	}

	newMarhsalled, err := yaml.Marshal(&types.Credentials{
		Active:   login,
		Keychain: credSet.Keychain,
	})
	if err != nil {
		return empty, err
	}

	fs.Replace(cp, newMarhsalled)

	return credMatch, nil
}
