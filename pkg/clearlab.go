package srklab

import (
        "os"
        "fmt"
        "os/exec"

        "sigs.k8s.io/kind/pkg/cmd"
        "github.com/cloud-native-everything/srklab/pkg/kind"
        "github.com/cloud-native-everything/srklab/pkg/dockerapi"
)

func Clear(configFile string) {

	kindVars := getVars(configFile)
	wg.Add(1)
	go delclab(kindVars.ClabTopology)
	delKind(*kindVars)
	wg.Wait()
	DockerNet.DockerNetworkDelete(kindVars.Network)

}

// Main is the kind main(), it will invoke Run(), if an error is returned
// it will then call os.Exit

func delKind(karray Kind_Data) {

	KindExec := func(Args *[]string) {
		if err := kindApp.Run(cmd.NewLogger(), cmd.StandardIOStreams(), *Args); err != nil {
			os.Exit(1)
		}
		wg.Done()
	}

	aux := []string{"delete", "clusters", "--all"}
	fmt.Println("--->", aux)
	wg.Add(1)
	go KindExec(&aux)

}


func delclab(mytopofile string) {

	// Create a new command to run "ls" with no arguments
	cmd := exec.Command("/usr/bin/containerlab", "destroy", "--topo", mytopofile)
	cmd.Stdout = &stdOut{}
	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	wg.Done()
}

