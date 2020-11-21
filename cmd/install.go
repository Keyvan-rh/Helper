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
	"github.com/robertsandoval/ocp4-helpernode/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)


var filename string
var helpme HelpMe



// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install creates a helpernode configuration",
	Long: `Install creates pulls images and sets up initial ~/.helpernodectl.yaml config file`,
	Run: func(cmd *cobra.Command, args []string) {
		readFile()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
//	installCmd.Flags().StringVarP(&filename, "filename", "f", "", "HelperNode file to create")
}


func readFile(){
	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &helpme)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	for name, image := range images {
		if(viper.GetBool("a_runtime." + name)) {
			utils.PullImage(image, DEFAULTTAG)
		}
	}
}