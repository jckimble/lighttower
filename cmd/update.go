package cmd

import (
	"fmt"
	"log"

	"github.com/jckimble/lighttower/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for update",
	Run:   startUpdate,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func startUpdate(cmd *cobra.Command, args []string) {
	var version string
	fmt.Sscanf(Version, "v%s", &version)
	u := &util.Updater{
		CurrentVersion: version,
		GithubOwner:    "jckimble",
		GithubRepo:     "lighttower",
	}
	available, err := u.CheckUpdateAvailable()
	if err != nil {
		log.Printf("Unable to check Update: %s\n", err)
	}

	if available != "" {
		log.Printf("Version %s available\n", available)
		err := u.Update()
		if err != nil {
			log.Printf("Unable to Update: %s\n", err)
		}
	} else {
		log.Println("LightTower is current Version")
	}
}
