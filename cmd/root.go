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
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command {
	Use:   "helpernodectl",
	Short: "Utility for the HelperNode",
	Long: `This cli utility is used to stop/start the HelperNode
on the host it's ran from. You need to provide a helpernode.yaml file
with information about your helper config. A simple example to start
your HelperNode is:

helpernodectl start --config=helpernode.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
//	verifyContainerRuntime()
//	verifyFirewallCommand()


	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helpernodectl.yaml)")
	//TODO Viper reads in ENV variables so not sure if there is a benefit for that way or this. Just adding this to research it.
	if len(os.Getenv("HELPERNODE_IMAGE_PREFIX")) > 0 {
		// Define images and their registry location based on the env var
		imageprefix := os.Getenv("HELPERNODE_IMAGE_PREFIX")
		images = map[string]string {
			"dns": imageprefix + "/helpernode/dns",
			"dhcp": imageprefix + "/helpernode/dhcp",
			"http": imageprefix + "/helpernode/http",
			"loadbalancer": imageprefix + "/helpernode/loadbalancer",
			"pxe": imageprefix + "/helpernode/pxe",
		}
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".helpernodectl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".helpernodectl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	setDefaults()
	reconcileImageList()


}

//This takes what was passed as --config and writes it to $HOME/.helpernodectl.yaml
func setDefaults(){
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".helpernodectl")
	viper.SetConfigType("yaml")
	//TODO figure out if we want to set defaults here

	//Touch the file in case it doesn't exist
	//TODO figure out a better way to do this
	emptyFile, err := os.Create(home + "/.helpernodectl.yaml")
	if err != nil {
		log.Fatal(err)
	}
	emptyFile.Close()

	err = viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config file")
		fmt.Println(err)
	} else {
		fmt.Println("Writing to:" + viper.ConfigFileUsed())
	}

}

func verifyContainerRuntime() {
	_, err := exec.LookPath("podman")
	if err != nil {
		fmt.Println("Podman not found, Please install")
		//TODO figure out if we really want to exit
		os.Exit(9)

	}

}

func verifyFirewallCommand() {

	_, err := exec.LookPath("firewall-cmd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error looking for firewall-cmd: %s\n", err)
		os.Exit(1)
	}
}


func reconcileImageList()  {
	disabledServices := viper.GetStringSlice("disabledServices")
	for name  := range disabledServices {
		delete(images, disabledServices[name])
	}
	//TODO Add code for pluginServices
}
