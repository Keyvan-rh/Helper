package cmd

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

//used by all commands to print output meaningfully
//logs as Info
//TODO maybe see if can bundle up a list of cmds to run
func runCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		logrus.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := scanner.Text()
		logrus.Info(m)
	}
	cmd.Wait()
}

//going to covert this to use the podman module in the future
func pullImage(image string, version string) {

	fmt.Println("Pulling: " + image + ":" + version)
	//TODO Need to write the output for the image pull
	cmd := exec.Command(containerRuntime, "pull", image+":"+version)
	runCmd(cmd)
}

//going to covert this to use the podman module in the future
//TODO clean up this to just take one string. build the string elsewhere
//TODO we need to adjust startImage to account for pluggable container that won't take the encoded file or need --net-host
func startImage(image string, version string, encodedyaml string, containername string) {
	logrus.Info("Starting helpernode-" + containername)
	cmd := exec.Command(containerRuntime, "run", "--rm", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--net=host", "--name=helpernode-"+containername, image+":"+version, "--label " + "helperpod_" + containername + "=" + VERSION )
	runCmd(cmd)

}

//going to covert this to use the podman module in the future
func stopImage(containername string) {

	logrus.Info("Stopping: helpernode-" + containername)
	//TODO check if image is runnign here rather than in start.go
	cmd := exec.Command(containerRuntime, "stop", "helpernode-"+containername)
	runCmd(cmd)

}
//check if an image is running. Return true if it is
//TODO see if we can do this with a --filter to get it to 1 result back. This implies building the iamge with some LABEL commands
func isImageRunning(containername string) bool {

	// output of all of all running containers
	out, err := exec.Command("podman", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		logrus.Debugf("Found %t running", out)
		name := string(out)
		if name == "helpernode" + containername {
			return true
		}
	}
/*
	// create a slice of string based on the output, trimming the newline first and splitting on "\n" (space)
	//TODO do we need  find() or do we tag the image and look for it directly
	s := strings.Split(strings.TrimSuffix(string(cmd), "\n"), "\n")
	_, found := find(s, containername)
	return found */
}

