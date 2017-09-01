package main

import (
	"./cmd"
	"fmt"
	"os"
)

var Version = "undefined"

func init() {
	cmd.Version = Version
}
func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
