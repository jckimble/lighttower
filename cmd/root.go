package cmd

import (
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version string
var cfgFile string
var RootCmd = &cobra.Command{
	Use:   "lighttower",
	Short: "tool to watch codeship builds and send dbus notifications",
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.lighttower.yaml", "Config file (default is $HOME/.lighttower.yaml)")
}

func initConfig() {
	cfgFile, err := homedir.Expand(cfgFile)
	if err != nil {
		log.Fatalf("Unable to find config: %s\n", err)
	}
	viper.SetConfigFile(cfgFile)
	viper.ReadInConfig()
}
