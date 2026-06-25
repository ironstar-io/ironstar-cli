package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InitCmd - `iron init`
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an Ironstar project",
	Long:  "Initializes an Ironstar project by creating configuration files in the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		err := services.InitializeIronstarProject()
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}

		color.Green("Successfully created the required Ironstar configuration files!")
		fmt.Println()
		color.Green("Run `iron subscription link [subscription-name]` to link this project to a remote Ironstar subscription")
	},
}

// InitIgnoreCmd - `iron init ignore`
var InitIgnoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "Migrate package excludes from config.yml to .ironstar/.ironignore",
	Long:  "Creates a .ironstar/.ironignore file from the existing config.yml package.exclude list (ported into their own group) plus sensible defaults, then removes the package block from config.yml. Package excludes are matched with full .gitignore semantics thereafter.",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := services.MigrateToIronignore()
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}

		color.Green("Created %s", res.IgnorePath)
		if len(res.Ported) > 0 {
			fmt.Printf("  • ported %d rule(s) from config.yml package.exclude into their own group\n", len(res.Ported))
		} else {
			fmt.Println("  • no package.exclude found in config.yml to port")
		}
		fmt.Println("  • added the standard Ironstar defaults")

		if res.RemovedPackage {
			color.Green("Removed the `package` block from %s", res.ConfigPath)
		} else {
			fmt.Printf("No `package` block found in %s; left unchanged\n", res.ConfigPath)
		}

		fmt.Println()
		color.Green("Package excludes are now governed by .ironstar/.ironignore (.gitignore syntax).")
		fmt.Println("Run `iron package --dry-run` to preview exactly what will be uploaded.")
	},
}
