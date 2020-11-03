package utils

import (
	"fmt"
	"os/exec"
)

var ImageName string
var ImageVersion string
//setting this up for future docker support possibly
var containerRuntime string = "podman"

//going to covert this to use the podman module in the future
func PullImage(image string, version string){

	fmt.Println("Pulling: " + image)
	//TODO Need to write the output for the image pull
	cmd, err := exec.Command(containerRuntime, "pull", image + ":" + version).Output()
	if err != nil {
		fmt.Println(cmd)
	}

}

//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string){

	fmt.Println("Running: " + image)
	//TODO Need to write the output for the image run
	cmd, err := exec.Command(containerRuntime, "run", "-d", "--env=HELPERPOD_CONFIG_YAML=" + encodedyaml, "--net=host", "--name=helpernode-" + containername, image + ":" + version).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(cmd)
	}

}
