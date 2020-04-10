package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/auth"
	"gitlab.com/ironstar-io/ironstar-cli/cmd/subscription"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd - `tok`
var rootCmd = cobra.Command{
	Use:   "ironstar",
	Short: "",
	Long:  "",
	Run:   run,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(VersionCmd)

	rootCmd.AddCommand(auth.AuthCmd)
	auth.AuthCmd.AddCommand(auth.LoginCmd)
	auth.AuthCmd.AddCommand(auth.ShowCmd)
	auth.AuthCmd.AddCommand(auth.SetProjectCredentialsCmd)
	auth.AuthCmd.AddCommand(auth.SetDefaultCredentials)

	rootCmd.AddCommand(subscription.SubscriptionCmd)
	subscription.SubscriptionCmd.AddCommand(subscription.ListCmd)
	subscription.SubscriptionCmd.AddCommand(subscription.LinkCmd)
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
	rootCmd.PersistentFlags().BoolP("force", "", false, "Forcefully skip destructive confirmation prompts")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "Auto-accept any non-destructive confirmation prompts")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode, command output is printed to the console")

	auth.LoginCmd.PersistentFlags().StringVarP(&auth.PasswordFlag, "password", "p", "", "Supply a password via the command line. Warning: Supplying the password via the command line is potentially insecure")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	if viper.GetBool("version") == true {
		fmt.Printf("v%s\n", version.Get().Version)
	} else {
		fmt.Printf("Ironstar v%s\n\n", version.Get().Version)
	}

	// çversion.SelfInstall(false)
}

func initConfig() {
	viper.BindPFlags(rootCmd.Flags())

	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}

	viper.ReadInConfig()
}
