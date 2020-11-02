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
		fmt.Println("here")
		fmt.Println(cmd)
	}

}