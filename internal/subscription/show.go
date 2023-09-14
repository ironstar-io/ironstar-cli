package subscription

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	if proj.Subscription == (types.Subscription{}) {
		return errors.New(errs.NoSubscriptionLinkErrorMsg)
	}

	if proj.Subscription.Alias == "" {
		return errors.New(errs.NoSubscriptionLinkErrorMsg)
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(proj)
		return nil
	}

	color.Green("Currently linked: ")
	fmt.Println(proj.Subscription.Alias + " (" + proj.Subscription.HashedID + ")")

	return nil
}
