package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Flags struct {
	DecryptFolder   string
	DecryptRemote   string
	EncryptRemote   string
	LocalFolder     string
	PlexDriveFolder string
	UnionFolder     string
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sbmt",
	Short: "Plexdrive, Rclone and UnionFS mount management made easier",
	Long:  RootLong,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sbmt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
