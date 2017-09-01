package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create config file with login info",
	Run:   startCreate,
	Args:  validateCreate,
}

var username, password string

func init() {
	createCmd.Flags().StringVarP(&username, "username", "u", "", "Codeship username/email")
	createCmd.Flags().StringVarP(&password, "password", "p", "", "Codeship password")
	RootCmd.AddCommand(createCmd)
}

func validateCreate(cmd *cobra.Command, args []string) error {
	if username == "" || password == "" {
		return fmt.Errorf("%s", "Username and Password Must be Provided")
	}
	return nil
}

func startCreate(cmd *cobra.Command, args []string) {
	yml := map[string]string{
		"auth": base64.StdEncoding.EncodeToString([]byte(username + ":" + password)),
	}
	data, _ := yaml.Marshal(yml)
	fmt.Printf("Writing to %s\n", cfgFile)
	cfgFile, _ := homedir.Expand(cfgFile)
	err := ioutil.WriteFile(cfgFile, data, 0700)
	if err != nil {
		fmt.Println(err)
	}
}
