package version

import (
	"fmt"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version/goos"

	"github.com/blang/semver"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Upgrade - Check if update available and auto-upgrade the user
func Upgrade(flg flags.Accumulator) error {
	// v := Get().Version
	// cv := strings.Replace(v, "v", "", -1)
	cs, err := semver.Parse("0.4.0")
	if err != nil {
		return errors.Wrap(err, "Unable to correctly parse the current version")
	}

	ls, err := GetLatest()
	if err != nil {
		return errors.Wrap(err, "Unable to reach GitHub to determine the latest version")
	}

	lv, _ := semver.Parse(strings.Replace(ls.TagName, "v", "", -1))
	if err != nil {
		return errors.Wrap(err, "Unable to correctly parse the latest version")
	}
	fmt.Println(cs)

	// Checks if the latest version is Greater Than the current version
	if lv.GT(cs) {
		confirmUpgrade := services.ConfirmationPrompt("This will upgrade your Ironstar CLI version to latest ("+lv.String()+").\n\nAre you sure?", "y", flg.AutoAccept)
		if !confirmUpgrade {
			fmt.Println("Exiting...")
			return nil
		}

		ip := GetInstallPath(lv.String())
		// Empty string if not installed, in which case, download and symlink
		if ip == "" {
			// Download latest release from GH
			_, err := DownloadAndInstall(lv.String())
			if err != nil {
				return errors.Wrap(err, "Ironstar CLI wasn't able to upgrade you to the latest version")
			}
		}

		// Latest version is installed, just not active. Create a Symlink to finish
		goos.ActivateSavedVersion(lv.String())

		fmt.Println()
		color.Green("Successfully upgraded the Ironstar CLI to the latest version (" + lv.String() + ")")

		return nil
	}

	color.Green("No updates available to the Ironstar CLI at this time. Version " + lv.String() + " is the latest available.")

	return nil
}
