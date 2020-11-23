package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Services []string `yaml:"disableservice"`
}
/*
//going to covert this to use the podman module in the future
func PullImage(image string, version string) {

	fmt.Println("Pulling: " + image)
	cmd, err := exec.Command(containerRuntime, "pull", image+":"+version).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

}
*/
// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
/*TODO need to remove this after testing
//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string) {

	fmt.Println("Running: " + image)
	cmd, err := exec.Command(containerRuntime, "run", "--rm", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--net=host", "--name=helpernode-"+containername, image+":"+version).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

}

//going to covert this to use the podman module in the future
func StopImage(containername string) {

	fmt.Println("Stopping: helpernode-" + containername)
	// First, stop container
	exec.Command(containerRuntime, "stop", "helpernode-"+containername).Output()
	// Then, rm the container so we can reuse the name afterwards
	exec.Command(containerRuntime, "rm", "--force", "helpernode-"+containername).Output()
}

 */

//check if an image is running. Return true if it is
//TODO see if we can do this with a --filter to get it to 1 result back. This implies building the iamge with some LABEL commands
func IsImageRunning(containername string) bool {

	// output of all of all running containers
	cmd, err := exec.Command("podman", "ps", "--format", "{{.Names}}").Output()

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on "\n" (space)
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), "\n")
	_, found := Find(s, containername)
	return found
}

// checking if service is running
func IsServiceRunning(servicename string) bool {
	// check if the service is active
	activestate, err := exec.Command("systemctl", "show", "-p", "ActiveState", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", activestate, err)
		os.Exit(53)
	}
	// return the status
	as := strings.TrimSuffix(strings.Split(string(activestate), "=")[1], "\n")
	return as == "active"
}

// checking if service is running
func IsServiceEnabled(servicename string) bool {
	// check if the service is active
	enabledstate, err := exec.Command("systemctl", "is-enabled", servicename).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", enabledstate, err)
		os.Exit(53)
	}
	// return the status
	es := strings.TrimSuffix(string(enabledstate), "\n")
	return es == "enabled"
}

// stopping service
func StopService(servicename string) {

	// stop the service only if it's running
	if IsServiceRunning(servicename) {
		fmt.Println("Stopping service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "stop", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// stopping service
func StartService(servicename string) {

	// start the service only if it isn't running
	if !IsServiceRunning(servicename) {
		fmt.Println("Starting service: " + servicename)
		//Start the service with systemd
		cmd, err := exec.Command("systemctl", "start", servicename).Output()
		// Check to see if the start was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// disable service
func DisableService(servicename string) {

	// Disable only if it needs to be
	if IsServiceEnabled(servicename) {
		fmt.Println("Disabling service: " + servicename)
		//Stop the service with systemd
		cmd, err := exec.Command("systemctl", "disable", servicename).Output()
		// Check to see if the stop was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// enable service
func EnableService(servicename string) {

	// Enable only if it needs to be
	if !IsServiceEnabled(servicename) {
		fmt.Println("Enabling service: " + servicename)
		//enable the service with systemd
		cmd, err := exec.Command("systemctl", "enable", servicename).Output()
		// Check to see if the enable was successful
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
			os.Exit(53)
		}
	}
}

// get current firewalld rules and return as a slice of string
func GetCurrentFirewallRules() []string {

	// get list of ports currently configured
	cmd, err := exec.Command("firewall-cmd", "--list-ports").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), " ")

	// get the list of services currenly configured
	scmd, err := exec.Command("firewall-cmd", "--list-services").Output()

	// check for error
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", scmd, err)
		os.Exit(253)
	}

	// create a slice of string based on the output, trimming the newline first and splitting on " " (space)
	svc := strings.Split(strings.TrimSuffix(string(scmd), "\n"), " ")

	// create a new array based on this new svc array. We will be converting service names to port output
	// simiar to what we got with: firewall-cmd --list--ports
	var ns = []string{}

	// range over the service, find out it's port and append it to the array we just created
	for _, v := range svc {
		lc, err := exec.Command("firewall-cmd", "--service", v, "--get-ports", "--permanent").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command %s: %s\n", lc, err)
			os.Exit(253)
		}
		nv := strings.TrimSuffix(string(lc), "\n")
		if strings.Contains(nv, " ") {
			ls := strings.Split(nv, " ")
			for _, l := range ls {
				ns = append(ns, l)
			}
		} else {
			ns = append(ns, nv)
		}
	}

	// append this new array of string into the original
	for _, v := range ns {
		s = append(s, v)
	}

	// Let's return this slice of string
	return s
}

func OpenPort(port string) {

	// Open Ports using the port number
	cmd, err := exec.Command("firewall-cmd", "--add-port", port, "--permanent", "-q").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running add-port command %s: %s\n", cmd, err)
		os.Exit(253)
	}

	// Reload the firewall to get the most up to date table
	rcmd, err := exec.Command("firewall-cmd", "--reload").Output()

	// check for error of command
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running reload command %s: %s\n", rcmd, err)
		os.Exit(253)
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

//This reconciles a list of images to start or stop
//defaults to all images unless specifically
func reconcileImageList(list []string) {

	disabledServices := viper.GetStringSlice("disabledServices")

	//use cases
	//all is implied so need to remove disabledServices
	if list[0]=="all"{
		//lets remove any disabled images
		for name  := range disabledServices {
			delete(images, disabledServices[name])
		}
	} else {
	// else we need to start with images and delete
	//  if images[name] == list[0] then do nothing else delete
		for _, name := range list{
			if _, exists := images[name]; !exists {
				delete(images, images[name])
			}
		}
	}
	//TODO
	//add plugable images
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
func validateArgs(args []string){
	imageCount := len(args)

	//if bare start command assume "all"
	if imageCount == 0 {
		imageList = []string{"all"}
	} else if imageCount == 1 {
		//parse image list
		imageList = strings.Split(args[0], ",")

		//TODO make sure plugable images is added to images var
		//Lets make sure its in our list of images (should include pluggable images)
		for _, name := range imageList {
			if _, exists := images[name]; exists {
				continue
			} else {
				logrus.Fatal("Listed service is not part of image list ")
			}

		}
	} else {
		logrus.Fatal("Wrong number of arguments passed. Must be comma separated list")
	}

}