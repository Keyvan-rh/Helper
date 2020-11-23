package cmd

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
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
func PullImage(image string, version string) {

	fmt.Println("Pulling: " + image + ":" + version)
	//TODO Need to write the output for the image pull
	cmd := exec.Command(containerRuntime, "pull", image+":"+version)
	runCmd(cmd)
}

//going to covert this to use the podman module in the future
//TODO clean up this to just take one string. build the string elsewhere
func StartImage(image string, version string, encodedyaml string, containername string) {
	//TODO check if image is runnign here rather than in start.go
	logrus.Info("Starting helpernode-" + containername)
	cmd := exec.Command(containerRuntime, "run", "--rm", "-d", "--env=HELPERPOD_CONFIG_YAML="+encodedyaml, "--net=host", "--name=helpernode-"+containername, image+":"+version)
	runCmd(cmd)

}

//going to covert this to use the podman module in the future
func StopImage(containername string) {

	logrus.Info("Stopping: helpernode-" + containername)
	//TODO check if image is runnign here rather than in start.go
	cmd := exec.Command(containerRuntime, "stop", "helpernode-"+containername)
	runCmd(cmd)

}

/*
//TODO  look at this function its from kind
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
