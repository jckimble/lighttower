package cmd

import (
	"errors"
	"log"

	"github.com/jckimble/lighttower/util"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start watching for builds",
	Run:   start,
	Args:  validate,
}

func init() {
	watchCmd.Flags().StringP("auth", "a", "", "Codeship authstring")
	viper.BindPFlag("auth", watchCmd.Flags().Lookup("auth"))

	viper.SetDefault("SuccessImage", "dialog-information")
	viper.SetDefault("SuccessMessage", "Build Successful")
	viper.SetDefault("SuccessSound", "")
	viper.SetDefault("ErrorImage", "dialog-error")
	viper.SetDefault("ErrorMessage", "Build Failed")
	viper.SetDefault("ErrorSound", "")
	RootCmd.AddCommand(watchCmd)
}

var validateMsg = "Codeship's Auth String must be set in config file or by flag"

func validate(cmd *cobra.Command, args []string) error {
	if viper.GetString("auth") == "" {
		return errors.New(validateMsg)
	}
	return nil
}
func notify(name, status string) {
	if status == "success" {
		sound, err := homedir.Expand(viper.GetString("SuccessSound"))
		if err != nil {
			sound = viper.GetString("SuccessSound")
		}
		err = util.SendNotify(viper.GetString("SuccessImage"), name, viper.GetString("SuccessMessage"), sound)
		if err != nil {
			log.Printf("Unable to send Notification: %s", err)
		}
	} else if status == "error" {
		sound, err := homedir.Expand(viper.GetString("ErrorSound"))
		if err != nil {
			sound = viper.GetString("ErrorSound")
		}
		err = util.SendNotify(viper.GetString("ErrorImage"), name, viper.GetString("ErrorMessage"), sound)
		if err != nil {
			log.Printf("Unable to send Notification: %s", err)
		}
	} else {
		log.Printf("Codeship returned status: %s\n", status)
	}
}
func start(cmd *cobra.Command, args []string) {
	codeship := util.CodeShip{
		AuthString: viper.GetString("auth"),
		Builds:     map[string]*util.Build{},
		CallBack:   notify,
	}
	err := codeship.GetToken()
	if err != nil {
		log.Fatalf("Unable To Get Projects: %s\n", err)
	}
	codeship.PollChanges()
}
