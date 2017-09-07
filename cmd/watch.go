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
	RootCmd.AddCommand(watchCmd)
}

var validateMsg = "Codeship's Auth String must be set in config file or by flag"

func validate(cmd *cobra.Command, args []string) error {
	if viper.GetString("auth") == "" {
		return errors.New(validateMsg)
	}
	return nil
}

func start(cmd *cobra.Command, args []string) {
	codeship := util.CodeShip{
		AuthString: viper.GetString("auth"),
		Builds:     map[string]*util.Build{},
	}
	err := codeship.GetToken()
	if err != nil {
		log.Fatalf("Unable To Get Projects: %s\n", err)
	}
	codeship.PollChanges()
}
