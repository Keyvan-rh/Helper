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
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helpernodectl",
	Short: "A tool to help with OCP installs",
	Long: `tool to help with OCP installs`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	fmt.Println("in init")

	verifyContainerRuntime()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helpernodectl.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("in initConfig")
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
	testConfigFileWrite()

}
func testConfigFileWrite(){
	home, err := homedir.Dir()

	fmt.Println("in testconfigfilewrite writing to " + home)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".helpernodectl")
	viper.SetConfigType("yaml")
	viper.SetDefault("ContentDir", "content")

	viper.WriteConfig()
}

func verifyContainerRuntime() {
	_, err := exec.LookPath("podman")
	if err != nil {
		fmt.Println("Podman not found, Please install")
		//TODO figure out if we really want to exit
	//	os.Exit(9)

	}

}