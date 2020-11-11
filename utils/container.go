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

	fmt.Println("Pulling: " + image + ":" + version)
	//TODO Need to write the output for the image pull
	cmd, err := exec.Command(containerRuntime, "pull", image + ":" + version).Output()
	if err != nil {
		fmt.Println(cmd)
		fmt.Println(err)
	}

}

//going to covert this to use the podman module in the future
func StartImage(image string, version string, encodedyaml string, containername string){

	/* TODO:
		- Need to write the output for the image run
		- Check if the image is already running
	*/
	_, err := exec.Command(containerRuntime, "run", "-d", "--env=HELPERPOD_CONFIG_YAML=" + encodedyaml, "--net=host", "--name=helpernode-" + containername, image + ":" + version).Output()
	if err != nil {
		fmt.Println(err)
	//		fmt.Println(cmd)
	}

}

//going to covert this to use the podman module in the future
func StopImage(containername string){

	fmt.Println("Stopping: helpernode-" + containername)
	/* TODO:
		- Need to write the output for the image run
		- Check if service is already stopped
	*/
	// First, stop container
	exec.Command(containerRuntime, "stop", "helpernode-" + containername).Output()
	// Then, rm the container so we can reuse the name afterwards
	rmcmd, err := exec.Command(containerRuntime, "rm", "--force", "helpernode-" + containername).Output()
	if err != nil {
		fmt.Println(err)
		fmt.Println(rmcmd)
	}

}

/*
//TODO  look at this function
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