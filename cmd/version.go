package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run:   getVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", Version)
}
