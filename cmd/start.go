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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// startCmd represents the start command
/* TODO need to remove once clean up is done
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts helpernode containers - must run as root",
	Long: "Starts helpernode containers - must run as root",
	Run: func(cmd *cobra.Command, args []string) {
		runContainers()
	},
}
 */

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts HelperNode containers based on the provided manifest.",
	Long: `This will start the containers needed for the HelperNode to run.
It will run the services depending on what manifest is passed.
Examples:
	helpernodectl start --config=helpernode.yaml
	
	cp helpernode.yaml ~/.helpernode.yaml
	helpernodectl start
This manifest should have all the information the services need to start
up successfully.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get any options passed
		skippreflight, _ := cmd.Flags().GetBool("skip-preflight")
		svc, _ := cmd.Flags().GetString("service")
		///
		if len(svc) > 0 {
			// check if the string passed matches a service
			if _, exists := images[svc]; exists {
				//if it matches; create a "single service" map and pass that to the stop function
				singleservicemap := map[string]string{svc: images[svc]}
				startImages(singleservicemap)
			} else {
				// If I didn't find it...tell them
				fmt.Println("Invalid service: " + svc)
				os.Exit(12)
			}
		} else {
			if skippreflight {
				fmt.Printf("Skipping Preflightchecks\n======================\n")
				startImages(images)
			} else {
				preflightCmd.Run(cmd, []string{})
				fmt.Printf("Starting Containers\n======================\n")
				startImages(images)
			}
		}
		///
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
	reconcileImageList()
	for name, image := range images {
		if IsImageRunning("helpernode-" + name) {
			fmt.Println("SKIPPING: Container helpernode-" + name + " already running.")
		} else {
			StartImage(image, "latest", getEncodedConfuration(), name)
		}
	}
}

func getEncodedConfuration() string {
	// Check to see if file exists
	logrus.Trace("Config file used: " + viper.ConfigFileUsed())
	var encoded string
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
		encoded = base64.StdEncoding.EncodeToString(content)
	}
	return encoded
}

