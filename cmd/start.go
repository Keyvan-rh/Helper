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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		validateArgs(args)
		return nil
	},
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
		skippreflight, _ := cmd.Flags().GetBool("skip-preflight")
		if skippreflight {
			logrus.Info("Skipping Preflightchecks\n======================\n")
		} else {
			preflightCmd.Run(cmd, []string{})
			logrus.Info("Starting Containers\n======================\n")
		}
		runContainers()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("skip-preflight", "", false, "Skips preflight checks and tries to start the containers")

}

func runContainers() {
	reconcileImageList(imageList)
	if logrus.GetLevel().String() == "debug" {
		for _, name := range imageList {
			logrus.Debug("Starting: " + name)
		}
	}
	for name, image := range images {
		if isImageRunning("helpernode-" + name) {
			logrus.Info("SKIPPING: Container helpernode-" + name + " already running.")
		} else {
			startImage(image, "latest", getEncodedConfuration(), name)
		}
	}
}
