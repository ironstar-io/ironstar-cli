package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/auth"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/backup"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/cache"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/deploy"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/env_vars"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/environment"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/logs"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/pkg"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/restore"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/subscription"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/sync"
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

	// `iron env-vars x` alias (hidden)
	rootCmd.AddCommand(env_vars.EnvVarsCmd)
	env_vars.EnvVarsCmd.AddCommand(env_vars.ListCmd)
	env_vars.EnvVarsCmd.AddCommand(env_vars.AddCmd)
	env_vars.EnvVarsCmd.AddCommand(env_vars.RemoveCmd)
	env_vars.EnvVarsCmd.AddCommand(env_vars.ModifyCmd)

	// `iron backup x`
	rootCmd.AddCommand(backup.BackupCmd)
	backup.BackupCmd.AddCommand(backup.NewCmd)
	backup.BackupCmd.AddCommand(backup.InfoCmd)
	backup.BackupCmd.AddCommand(backup.ListCmd)
	backup.BackupCmd.AddCommand(backup.DeleteCmd)
	backup.BackupCmd.AddCommand(backup.DownloadCmd)

	// `iron restore x`
	rootCmd.AddCommand(restore.RestoreCmd)
	restore.RestoreCmd.AddCommand(restore.NewCmd)
	restore.RestoreCmd.AddCommand(restore.InfoCmd)
	restore.RestoreCmd.AddCommand(restore.ListCmd)

	// `iron sync x`
	rootCmd.AddCommand(sync.SyncCmd)
	sync.SyncCmd.AddCommand(sync.NewCmd)
	sync.SyncCmd.AddCommand(sync.InfoCmd)
	sync.SyncCmd.AddCommand(sync.ListCmd)

	// `iron logs x`
	rootCmd.AddCommand(logs.LogsCmd)

	// `iron cache x`
	rootCmd.AddCommand(cache.CacheCmd)
	cache.CacheCmd.AddCommand(cache.InvalidationCmd)
	cache.InvalidationCmd.AddCommand(cache.ListInvalidationsCmd)
	cache.InvalidationCmd.AddCommand(cache.ShowInvalidationCmd)
	cache.InvalidationCmd.AddCommand(cache.CreateCmd)

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
	rootCmd.PersistentFlags().StringArrayVarP(&flags.Acc.Component, "component", "c", []string{}, "Supply an array of components to backup/restore/sync")

	LoginCmd.PersistentFlags().StringVarP(&flags.Acc.Password, "password", "p", "", "Supply a password via the command line. Warning: Supplying the password via the command line is potentially insecure")

	deploy.DeployCmd.PersistentFlags().BoolVarP(&flags.Acc.SkipHooks, "skip-hooks", "", false, "Skip running the hooks defined in .ironstar/config.yml on this release")
	deploy.DeployCmd.PersistentFlags().BoolVarP(&flags.Acc.PreventRollback, "prevent-rollback", "", false, "Prevent automated rollback of a deployment in case of failure.")

	backup.BackupCmd.PersistentFlags().BoolVarP(&flags.Acc.LockTables, "lock-tables", "", false, "Pass the `--lock-tables` flag to mysqldump")
	backup.NewCmd.PersistentFlags().BoolVarP(&flags.Acc.LockTables, "lock-tables", "", false, "Pass the `--lock-tables` flag to mysqldump")
	backup.DownloadCmd.PersistentFlags().StringVarP(&flags.Acc.SavePath, "save-path", "", "", "Select a custom save path for your backup download. Default: ./.ironstar/backups/{subscription}/{environment}")
	backup.DownloadCmd.PersistentFlags().StringVarP(&flags.Acc.Backup, "backup", "", "", "Select the backup to be the base for download")

	// TODO
	logs.LogsCmd.PersistentFlags().StringVarP(&flags.Acc.Backup, "backup", "", "", "Select the backup to be the base for download")

	restore.RestoreCmd.PersistentFlags().StringVarP(&flags.Acc.Strategy, "strategy", "", "", "Provide the strategy for a restore")
	restore.NewCmd.PersistentFlags().StringVarP(&flags.Acc.Strategy, "strategy", "", "", "Provide the strategy for a restore")
	restore.RestoreCmd.PersistentFlags().StringVarP(&flags.Acc.Backup, "backup", "", "", "The source backup identifier to restore from")
	restore.NewCmd.PersistentFlags().StringVarP(&flags.Acc.Backup, "backup", "", "", "The source backup identifier to restore from")

	sync.SyncCmd.PersistentFlags().StringVarP(&flags.Acc.SrcEnvironment, "src", "", "", "Identifies the source environment to copy from")
	sync.NewCmd.PersistentFlags().StringVarP(&flags.Acc.SrcEnvironment, "src", "", "", "Identifies the source environment to copy from")
	sync.SyncCmd.PersistentFlags().StringVarP(&flags.Acc.DestEnvironment, "dest", "", "", "Identifies the destination environment to copy to")
	sync.NewCmd.PersistentFlags().StringVarP(&flags.Acc.DestEnvironment, "dest", "", "", "Identifies the destination environment to copy to")
	sync.SyncCmd.PersistentFlags().StringVarP(&flags.Acc.SrcEnvironment, "src-env", "", "", "Identifies the source environment to copy from")
	sync.NewCmd.PersistentFlags().StringVarP(&flags.Acc.SrcEnvironment, "src-env", "", "", "Identifies the source environment to copy from")
	sync.SyncCmd.PersistentFlags().StringVarP(&flags.Acc.DestEnvironment, "dest-env", "", "", "Identifies the destination environment to copy to")
	sync.NewCmd.PersistentFlags().StringVarP(&flags.Acc.DestEnvironment, "dest-env", "", "", "Identifies the destination environment to copy to")
	sync.SyncCmd.PersistentFlags().BoolVarP(&flags.Acc.UseLatestBackup, "use-latest-backup", "", false, "Use this flag to instruct this operation to use the latest full scheduled backup from the --source-env as the source, this will prevent a new backup from being taken and will significantly improve the sync time.")
	sync.NewCmd.PersistentFlags().BoolVarP(&flags.Acc.UseLatestBackup, "use-latest-backup", "", false, "Use this flag to instruct this operation to use the latest full scheduled backup from the --source-env as the source, this will prevent a new backup from being taken and will significantly improve the sync time.")
	sync.SyncCmd.PersistentFlags().StringVarP(&flags.Acc.Strategy, "strategy", "", "", "Provide the strategy for the restore section of the sync")
	sync.NewCmd.PersistentFlags().StringVarP(&flags.Acc.Strategy, "strategy", "", "", "Provide the strategy for the restore section of the sync")

	env_vars.EnvVarsCmd.PersistentFlags().StringVarP(&flags.Acc.Key, "key", "", "", "The environment variable key")
	env_vars.EnvVarsCmd.PersistentFlags().StringVarP(&flags.Acc.Value, "value", "", "", "The environment variable value")
	env_vars.EnvVarsCmd.PersistentFlags().StringVarP(&flags.Acc.VarType, "var-type", "", "PROTECTED", "Either PROTECTED or VISIBLE. Values VISIBLE are encrypted in the Ironstar database, but visible in plaintext to authenticated users of the Ironstar API/UI. Values PROTECTED are encrypted in the Ironstar database and never visible in plaintext to the API/UI. Default PROTECTED")

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
