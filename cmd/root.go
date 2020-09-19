package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/auth"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/backup"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/deploy"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/environment"
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
	rootCmd.AddCommand(UpgradeCmd)
	rootCmd.AddCommand(LoginCmd)
	rootCmd.AddCommand(LogoutCmd)
	rootCmd.AddCommand(InitCmd)

	// `iron auth x`
	rootCmd.AddCommand(auth.AuthCmd)
	auth.AuthCmd.AddCommand(auth.ShowCmd)
	auth.AuthCmd.AddCommand(auth.SetActiveCmd)
	auth.AuthCmd.AddCommand(auth.MFACmd)
	auth.MFACmd.AddCommand(auth.MFADisableCmd)
	auth.MFACmd.AddCommand(auth.MFAEnableCmd)

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

	// `iron environment x`
	rootCmd.AddCommand(environment.EnvironmentCmd)
	environment.EnvironmentCmd.AddCommand(environment.ListCmd)
	environment.EnvironmentCmd.AddCommand(environment.DisableRestoreCmd)
	environment.EnvironmentCmd.AddCommand(environment.EnableRestoreCmd)

	// `iron env x` alias (hidden)
	rootCmd.AddCommand(environment.EnvCmd)
	environment.EnvCmd.AddCommand(environment.ListCmd)
	environment.EnvCmd.AddCommand(environment.DisableRestoreCmd)
	environment.EnvCmd.AddCommand(environment.EnableRestoreCmd)

	// `iron backup x`
	rootCmd.AddCommand(backup.BackupCmd)
	backup.BackupCmd.AddCommand(backup.NewCmd)
	backup.BackupCmd.AddCommand(backup.InfoCmd)

	// `iron package x`
	rootCmd.AddCommand(pkg.PackageCmd)
	pkg.PackageCmd.AddCommand(pkg.ListCmd)
	pkg.PackageCmd.AddCommand(pkg.UpdateRefCmd)

	// `iron pkg x` alias (hidden)
	rootCmd.AddCommand(pkg.PkgCmd)
	pkg.PkgCmd.AddCommand(pkg.ListCmd)
	pkg.PkgCmd.AddCommand(pkg.UpdateRefCmd)

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
	// rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode, command output is printed to the console")
	rootCmd.PersistentFlags().BoolVarP(&flags.Acc.AutoAccept, "yes", "y", false, "Auto-accept any non-destructive confirmation prompts")

	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Login, "login", "l", "", "Force use of a specified logins' credentials")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Output, "output", "o", "", "Use a certain output type. Not applicable on all commands.")

	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Subscription, "subscription", "s", "", "Use or filter by a specified subscription. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Environment, "environment", "e", "", "Use or filter by  specified environment. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Environment, "env", "", "", "Use or filter by  specified environment. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Package, "package", "", "", "Use or filter by  specified package. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Deploy, "deploy", "d", "", "Use or filter by  specified deployment. Not applicable on all commands.")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Exclude, "exclude", "", "", "A comma separated list of files/directories to exclude during packaging")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Name, "name", "n", "", "Supply a name, not applicable for all command")
	rootCmd.PersistentFlags().StringVarP(&flags.Acc.Type, "type", "t", "", "Supply a type, not applicable for all command")
	rootCmd.PersistentFlags().StringArrayVarP(&flags.Acc.Component, "component", "c", []string{""}, "Supply an array of components to backup/restore/sync")

	LoginCmd.PersistentFlags().StringVarP(&flags.Acc.Password, "password", "p", "", "Supply a password via the command line. Warning: Supplying the password via the command line is potentially insecure")

	backup.BackupCmd.PersistentFlags().StringVarP(&flags.Acc.Retention, "retention", "r", "", "Provide the retention period for a backup")
	backup.NewCmd.PersistentFlags().StringVarP(&flags.Acc.Retention, "retention", "r", "", "Provide the retention period for a backup")

	pkg.PkgCmd.PersistentFlags().StringVarP(&flags.Acc.Ref, "ref", "", "", "A user defined reference used for being able to easily identify the package. This could be a git commit SHA, UUID, or tag of your choice. It is not mandatory.")
	pkg.PackageCmd.PersistentFlags().StringVarP(&flags.Acc.Ref, "ref", "", "", "A user defined reference used for being able to easily identify the package. This could be a git commit SHA, UUID, or tag of your choice. It is not mandatory.")
	pkg.UpdateRefCmd.PersistentFlags().StringVarP(&flags.Acc.Ref, "ref", "", "", "A user defined reference used for being able to easily identify the package. This could be a git commit SHA, UUID, or tag of your choice. It is not mandatory.")
	deploy.DeployCmd.PersistentFlags().StringVarP(&flags.Acc.Ref, "ref", "", "", "A user defined reference used for being able to easily identify the package. This could be a git commit SHA, UUID, or tag of your choice. It is not mandatory.")

	deploy.DeployCmd.PersistentFlags().BoolVarP(&flags.Acc.ApproveProdDeploy, "approve-prod-deploy", "", false, "Allow deployments to production without a warning prompt.")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	if viper.GetBool("version") {
		fmt.Printf("%s\n", version.Get().Version)
	}
	fmt.Printf("Ironstar %s\n\n", version.Get().Version)
}

func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}

	viper.ReadInConfig()
}
