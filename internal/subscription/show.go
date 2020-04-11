package subscription

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Show(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	confPath := filepath.Join(wd, ".ironstar", "config.yml")

	exists := fs.CheckExists(confPath)
	if !exists {
		return errors.New("No Ironstar configuration found in this directory. Have you run `iron init`")
	}

	proj, err := services.GetProjectData()
	if err != nil {
		return err
	}

	if proj.Subscription.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Currently linked: ")
	fmt.Println(proj.Subscription.Alias)

	return nil
}
