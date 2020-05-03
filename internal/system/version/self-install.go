package version

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version/goos"

	"github.com/blang/semver"
)

// SelfInstall - Check state then install runing Ironstar CLI binary to PATH
func SelfInstall(forceInstall bool) {
	if forceInstall == true {
		// Ironstar CLI not in PATH, confirm install of current bin
		confirmUpgrade := services.ConfirmationPrompt("This command will install Ironstar CLI "+Get().Version+". Would you like to continue?", "y")
		if confirmUpgrade == false {
			fmt.Println("Exiting...")
			return
		}

		installRunningBin()

		return
	}

	_, err := exec.LookPath("tok")
	if err != nil {
		// Ironstar CLI not in PATH, confirm install of current bin
		confirmUpgrade := services.ConfirmationPrompt("It looks like this is your first time running Ironstar CLI.\n\nWould you like to install it now", "y")
		if confirmUpgrade == false {
			fmt.Println("Exiting...")
			return
		}

		installRunningBin()

		return
	}

	// Ironstar CLI already in PATH, display help message and exit
	// TODO
	fmt.Println("For help with Ironstar CLI run `tok help` or take a look at our documentation at https://docs.ironstar-cli.io")
}

// installRunningBin - Install runing Ironstar CLI binary to the global install path
func installRunningBin() {
	v := Get().Version
	cv := strings.Replace(v, "v", "", 0)
	cs, err := semver.Parse(cv)
	if err != nil {
		fmt.Println("Ironstar CLI was unable to correctly parse the current version: ", err.Error())
		os.Exit(1)
	}
	bv := cs.String()

	fmt.Println("Checking for an existing installation of Ironstar CLI version " + bv)
	ip := GetInstallPath(bv)
	// Empty string if not installed, in which case, save it
	if ip == "" {
		fmt.Println("This Ironstar CLI version (" + bv + ") is not installed")
		p, err := goos.SaveTokBinary(bv)
		if err != nil {
			fmt.Println("Ironstar CLI wasn't able to install this version correctly: ", err.Error())
			os.Exit(1)

		}

		ip = p
	}

	// This running instance is saved, now we activate it as the default 'installed' version
	fmt.Println("Making Ironstar CLI version [" + bv + "] the default version")
	goos.ActivateSavedVersion(bv)

	fmt.Println()
	fmt.Println("Success! Ironstar CLI version " + bv + " should now be avaliable as 'tok'")
}
