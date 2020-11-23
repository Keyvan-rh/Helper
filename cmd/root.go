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
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//this is really only here for testing the code on a mac
	//helpernodectl is not currently supported on macOS
	if runtime.GOOS != "darwin" {
		verifyContainerRuntime()
		verifyFirewallCommand()
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helpernodectl.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")

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


	createImageList()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Trace("Using config file:", viper.ConfigFileUsed())
	}
	setDefaults()
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
func createImageList(){
	viper.SetEnvPrefix("helpernode")
	viper.BindEnv("image_prefix")
	if viper.GetString("image_prefix") == "" {
		logrus.Debug("HELPERNODE_IMAGE_PREFIX not found. using quay.io")
		viper.Set("image_prefix", "quay.io")
	} else {
		logrus.Debug("here")
	}

	viper.AutomaticEnv() // read in environment variables that match

	registry = viper.GetString("image_prefix")

	for _, name := range coreImageNames{
		images[name]=registry + "/" + repository + "/"  + name
	}
	//TODO Add pluggable images here

	//Just some logic to print if in debug
	if logrus.GetLevel().String() == "debug" {
		logrus.Debug("Using registry : " + registry)
		for name,image := range images {
			logrus.Debug(name + ":" + image)
		}
	}
}