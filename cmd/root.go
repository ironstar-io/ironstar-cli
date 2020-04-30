package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/auth"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/deploy"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/pkg"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/subscription"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd - `iron`
var rootCmd = cobra.Command{
	Use:   "ironstar",
	Short: "",
	Long:  "",
	Run:   run,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(LoginCmd)
	rootCmd.AddCommand(InitCmd)

	// `iron auth x`
	rootCmd.AddCommand(auth.AuthCmd)
	auth.AuthCmd.AddCommand(auth.ShowCmd)
	auth.AuthCmd.AddCommand(auth.SetActiveCmd)
	auth.AuthCmd.AddCommand(auth.MFAEnableCmd)
	auth.AuthCmd.AddCommand(auth.MFADisableCmd)

	// `iron subscription x`
	rootCmd.AddCommand(subscription.SubscriptionCmd)
	subscription.SubscriptionCmd.AddCommand(subscription.ListCmd)
	subscription.SubscriptionCmd.AddCommand(subscription.LinkCmd)
	subscription.SubscriptionCmd.AddCommand(subscription.ShowCmd)

	// `iron sub x` alias (hidden)
	rootCmd.AddCommand(subscription.SubCmd)
	subscription.SubCmd.AddCommand(subscription.ListCmd)
	subscription.SubCmd.AddCommand(subscription.LinkCmd)
	subscription.SubCmd.AddCommand(subscription.ShowCmd)

	// `iron package x`
	rootCmd.AddCommand(pkg.PackageCmd)
	pkg.PackageCmd.AddCommand(pkg.ListCmd)

	// `iron pkg x` alias (hidden)
	rootCmd.AddCommand(pkg.PkgCmd)
	pkg.PkgCmd.AddCommand(pkg.ListCmd)

	// `iron deploy x`
	rootCmd.AddCommand(deploy.DeployCmd)
	deploy.DeployCmd.AddCommand(deploy.ListCmd)
	deploy.DeployCmd.AddCommand(deploy.StatusCmd)
}

// Execute - Root executable
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// RootCmd will setup and return the root command
func RootCmd() *cobra.Command {
	// rootCmd.PersistentFlags().BoolP("force", "", false, "Forcefully skip destructive confirmation prompts")
	// rootCmd.PersistentFlags().BoolP("yes", "y", false, "Auto-accept any non-destructive confirmation prompts")
	// rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode, command output is printed to the console")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Login, "login", "l", "", "Force use of a specified logins' credentials")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Output, "output", "o", "", "Use a certain output type. Not applicable on all commands.")

	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Subscription, "subscription", "s", "", "Use or filter by a specified subscription. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Environment, "environment", "e", "", "Use or filter by  specified environment. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Package, "package", "", "", "Use or filter by  specified package. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Deploy, "deploy", "d", "", "Use or filter by  specified deployment. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Exclude, "exclude", "", "", "A comma separated list of files/directories to exclude during packaging")

	LoginCmd.PersistentFlags().StringVarP(&flags.Acc.Password, "password", "p", "", "Supply a password via the command line. Warning: Supplying the password via the command line is potentially insecure")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	if viper.GetBool("version") == true {
		fmt.Printf("v%s\n", version.Get().Version)
	}
	fmt.Printf("Ironstar v%s\n\n", version.Get().Version)

	// version.SelfInstall(false)
}

func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}

	viper.ReadInConfig()
}
