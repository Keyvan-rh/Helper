package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
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
	setUpLogging()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	if runtime.GOOS != "darwin" {
		verifyContainerRuntime()
		verifyFirewallCommand()
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helpernodectl.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")


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
	setUpLogging()
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
		logrus.Trace("Using config file:", viper.ConfigFileUsed())
	}


	setDefaults()
	//TODO need to move this to start/stop or actions that need the image list
	reconcileImageList()


}
func setUpLogging() {
    logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log-level"))
	}

	logrus.SetLevel(level)

}

//This takes what was passed as --config and writes it to $HOME/.helpernodectl.yaml
func setDefaults(){
//	fmt.Println("root.go:setDefaults() called from root.go:initConfig()")
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
		logrus.Error("Error writing config file")
	} else {
		logrus.Trace("Writing to:" + viper.ConfigFileUsed())
	}

}
