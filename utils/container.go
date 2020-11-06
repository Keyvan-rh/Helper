package utils

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
)

var ImageName string
var ImageVersion string

//setting this up for future docker support possibly
var containerRuntime string = "podman"

//going to covert this to use the podman module in the future
func PullImage(image string, version string) {

	fmt.Println("Pulling: " + image + version)
	//TODO Need to write the output for the image pull
	cmd, err := exec.Command(containerRuntime, "pull", image+":"+version).Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(err)
	}

}

//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string) {

	fmt.Println("Running: " + image)
	/* TODO:
	- Need to write the output for the image run
	- Check if the image is already running
	*/
<<<<<<< HEAD
	cmd, err := exec.Command(containerRuntime, "run", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--net=host", "--name=helpernode-"+containername, image+":"+version).Output()
=======
	addMounts()
	cmd, err := exec.Command(containerRuntime, "run", "-d", "--env=HELPERPOD_CONFIG_YAML=" + encodedyaml, "--net=host", "--name=helpernode-" + containername, image + ":" + version).Output()
>>>>>>> 1435f739d21a50bef12c689a5d30693eb2744547
	if err != nil {
		fmt.Println(err)
		fmt.Println(cmd)
	}
	removeMounts()

}

func addMounts(){
	//copy -f file from create to $HOME/.config/contain
	//check if mounts.conf exists
}



func removeMounts(){
	//zero out mounts file or rather
}












//going to covert this to use the podman module in the future
func StopImage(containername string) {

	fmt.Println("Stopping: helpernode-" + containername)
	/* TODO:
	- Need to write the output for the image run
	- Check if service is already stopped
	*/
	// First, stop container
	exec.Command(containerRuntime, "stop", "helpernode-"+containername).Output()
	// Then, rm the container so we can reuse the name afterwards
	rmcmd, err := exec.Command(containerRuntime, "rm", "--force", "helpernode-"+containername).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(rmcmd)
	}

}

/*
// pull pulls an image, retrying up to retries times
func pull(logger log.Logger, image string, retries int) error {
	logger.V(1).Infof("Pulling image: %s ...", image)
	err := exec.Command("podman", "pull", image).Run()
	// retry pulling up to retries times if necessary
	if err != nil {
		for i := 0; i < retries; i++ {
			time.Sleep(time.Second * time.Duration(i+1))
			logger.V(1).Infof("Trying again to pull image: %q ... %v", image, err)
			// TODO(bentheelder): add some backoff / sleep?
			err = exec.Command("podman", "pull", image).Run()
			if err == nil {
				break
			}
		}
	}
	return errors.Wrapf(err, "failed to pull image %q", image)
}
*/
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
