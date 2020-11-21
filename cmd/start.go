/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/robertsandoval/ocp4-helpernode/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts helpernode containers - must run as root",
	Long: "Starts helpernode containers - must run as root",
	Run: func(cmd *cobra.Command, args []string) {
		runContainers()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("skip-preflight", "", false, "Skips preflight checks and tries to start the containers")
	startCmd.Flags().String("service", "", "start a service/container (preflight NOT performed). Valid names: dns, dhcp, http, loadbalancer, pxe")

	//TODO lets move the file read from ~/.helpernodectl.yaml or from --config here and remove the read from runContainers()
	// or better yet move it into root so as things get added its globally available and there is functionality there already
}

func runContainers() {
	// Check to see if file exists
	fmt.Println("config file used: " + viper.ConfigFileUsed())
	configurationFile:=viper.ConfigFileUsed()
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		fmt.Println("File " + configurationFile + " does not exist")
	} else {
		// Open file on disk
		f, _ := os.Open(configurationFile)
		// Read file into a byte slice
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)
		//Encode to base64
		encoded := base64.StdEncoding.EncodeToString(content)
		// run the containers using the encoding
		for name, image := range images {
			if IsImageRunning("helpernode-" + name) {
				fmt.Println("SKIPPING: Container helpernode-" + name + " already running.")
			} else {
				utils.StartImage(image, "latest", encoded, name)
			}
		}
	}
}

