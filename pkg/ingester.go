package srklab

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ExecKubeApps(myapps Resources_Data, mykubeconf string, myns string) {
	// Define arguments for the kubectl command
	kubeconfigfile := "--kubeconfig=" + mykubeconf

	args := []string{"apply", "-f", myapps.App, kubeconfigfile}

	// Create a command object with the kubectl command and its arguments
	cmd := exec.Command("kubectl", args...)

	// Run the command and capture its output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	// Check for errors
	if err != nil {
		fmt.Printf("Error running kubectl command: %v\n", err)
		fmt.Printf("Command output for %s: %s\n", mykubeconf, stderr.String())
	} else {
		fmt.Printf("Command output for %s: %s\n", mykubeconf, stdout.String())
	}

	if strings.Contains(myapps.App, "multus") {
		fmt.Printf("INFO: Waiting for Multus to be in running state at %s\n", mykubeconf)
		time.Sleep(70 * time.Second)
	}

}
