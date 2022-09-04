package version

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/ironstar-io/ironstar-cli/internal/services/github"
	"github.com/ironstar-io/ironstar-cli/internal/system/version/goos"
)

// GetInstallPath - Check if Ironstar CLI version is installed or not
func GetInstallPath(version string) string {
	return goos.GetInstallPath(version)
}

// DownloadAndInstall - Download and install a selected Ironstar CLI version
func DownloadAndInstall(version string) (string, error) {
	releaseTag, err := GetReleaseTagFromVersion(version)
	if err != nil {
		return "", err
	}

	return goos.DownloadCLIBinary(releaseTag)
}

// GetReleaseTagFromVersion returns a githab-ready release tag from a Ironstar CLI version string
func GetReleaseTagFromVersion(version string) (string, error) {
	var empty string
	// Check version exists in GH
	// Get the URL for the version
	ghr := []github.ReleaseBody{}
	res, err := github.GetAllReleases()
	if err != nil {
		return empty, errors.Wrap(err, "Unexpected error retrieving list of available Ironstar CLI releases:")
	}

	err = json.Unmarshal(res.Body, &ghr)
	if err != nil {
		return empty, errors.Wrap(err, "Unexpected error assembling list of available Ironstar CLI releases")
	}

	for _, r := range ghr {
		if r.TagName == version {
			if r.Draft {
				fmt.Println("\nWarning: The selected version is a draft and may not work as intended")
			}

			if r.Prerelease {
				fmt.Println("\nWarning: The selected version is a prerelease and may not work as intended")
			}

			return r.TagName, nil
		}
	}

	return empty, errors.New("There is no release information available for version [" + version + "]")
}
