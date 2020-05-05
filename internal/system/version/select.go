package version

import (
	"fmt"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version/goos"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

const minimumCLIVersion = "0.4.1"

// Select - Change the users' Ironstar CLI version to their selection
func Select(selection string, flg flags.Accumulator) error {
	v := Get().Version
	cv := strings.Replace(v, "v", "", -1)
	cs, _ := semver.Parse(cv)

	cutSelection := strings.Replace(selection, "v", "", -1)
	sv, err := semver.Parse(cutSelection)
	if err != nil {
		fmt.Println("Invalid semver selection supplied. Exiting...")
		return nil
	}

	selectedVersion := "v" + sv.String()

	// Checks if the current version is Equal to the selected version and exits if so
	if sv.EQ(cs) {
		fmt.Println("Selected version (" + selectedVersion + ") is the same as the currently active version. Exiting...")
		return nil
	}

	mv, _ := semver.Parse(minimumCLIVersion)
	// Checks if the selected version is Lesser Than the minimum version
	if sv.LT(mv) {
		fmt.Println("Selected version (" + selectedVersion + ") is less than the minimum allowed version (" + minimumCLIVersion + "). Exiting...")
		return nil
	}

	confirmChange := services.ConfirmationPrompt("This will change your Ironstar CLI version to "+selectedVersion+".\n\nAre you sure?", "y", flg.AutoAccept)
	if !confirmChange {
		fmt.Println("Exiting...")
		return nil
	}

	ip := GetInstallPath(selectedVersion)
	// Empty string if not installed, in which case, download and symlink
	if ip == "" {
		// Download & install selected release from GH
		_, err := DownloadAndInstall(selectedVersion)
		if err != nil {
			return errors.Wrap(err, "Ironstar CLI was unable to upgrade you to the selected version")
		}
	}

	// This running instance is saved, now we activate it as the default 'installed' version
	err = goos.ActivateSavedVersion(selectedVersion)
	if err != nil {
		return errors.Wrap(err, "Ironstar CLI was unable to change you to the selected version")
	}

	fmt.Println("Successfully changed Ironstar CLI to version " + selectedVersion)

	return nil
}
