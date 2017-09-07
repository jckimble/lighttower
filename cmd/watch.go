package cmd

import (
	"errors"
	"log"

	"github.com/jckimble/lighttower/util"

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
	viper.SetDefault("ErrorImage", "dialog-error")
	viper.SetDefault("ErrorMessage", "Build Failed")
	viper.SetDefault("Debug", false)
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
	if viper.GetBool("Debug") {
		log.Printf("Codeship returned status: %s\n", status)
	}
	if status == "success" {
		err := util.SendNotify(viper.GetString("SuccessImage"), name, viper.GetString("SuccessMessage"))
		if err != nil {
			log.Printf("Unable to send Notification: %s", err)
		}
	} else if status == "error" {
		err := util.SendNotify(viper.GetString("ErrorImage"), name, viper.GetString("ErrorMessage"))
		if err != nil {
			log.Printf("Unable to send Notification: %s", err)
		}
	} else if status != "testing" && !viper.GetBool("Debug") {
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
