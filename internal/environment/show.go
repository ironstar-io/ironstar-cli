package environment

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Show(args []string, flg flags.Accumulator) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if !exists {
		return errors.New("No Ironstar configuration found in this directory. Have you run `iron init`")
	}

	proj, err := services.GetProjectData(flg.AutoAccept)
	if err != nil {
		return err
	}

	if proj.Environment == (types.Environment{}) {
		return errors.New(errs.NoEnvironmentLinkErrorMsg)
	}

	if proj.Environment.Name == "" {
		return errors.New(errs.NoEnvironmentLinkErrorMsg)
	}

	class := "Non-Production"
	if proj.Environment.Class == "cw" {
		class = "Production"
	}

	color.Green("Currently linked environment for subscription '" + proj.Subscription.Alias + "':")

	fmt.Println()
	fmt.Println(proj.Environment.Name + " (" + class + ")")
	if proj.Environment.Class == "cw" {
		fmt.Println()
		color.Yellow("Caution! Your current context is a production environment!")
	}

	return nil
}
