package srklab

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	DockerNet "github.com/cloud-native-everything/srklab/pkg/dockerapi"
	kindApp "github.com/cloud-native-everything/srklab/pkg/kind"
	"sigs.k8s.io/kind/pkg/cmd"
)

func Start(configFile string) {

	kindVars := getVars(configFile)
	DockerNet.DockerNetworkCreate(kindVars.Network, kindVars.Prefix)
	wg.Add(1)
	go clab(kindVars.ClabTopology)
	runKind(*kindVars)
	wg.Wait()
	loadKind(*kindVars)
	wg.Wait()
	cklinks(kindVars.Links)
	kindsetInf(kindVars.Links)
	for i := 0; i < len(kindVars.Cluster); i++ {
		wg.Add(1)
		go execApps(kindVars.Cluster[i].Resources, kindVars.Cluster[i].Kubeconfig, "default")
	}
	wg.Wait()

}

// Main is the kind main(), it will invoke Run(), if an error is returned
// it will then call os.Exit

func runKind(karray Kind_Data) {

	KindExec := func(Args *[]string) {
		if err := kindApp.Run(cmd.NewLogger(), cmd.StandardIOStreams(), *Args); err != nil {
			os.Exit(1)
		}
		wg.Done()
	}

	for i := 0; i < len(karray.Cluster); i++ {
		aux := []string{"create", "cluster", "--name", karray.Cluster[i].Name}
		aux = append(aux, "--image", karray.Cluster[i].Image)
		aux = append(aux, "--config", karray.Cluster[i].Config)
		aux = append(aux, "--kubeconfig", karray.Cluster[i].Kubeconfig)
		fmt.Println("--->", aux)
		wg.Add(1)
		go KindExec(&aux)

	}
}

func loadKind(karray Kind_Data) {

	KindExec := func(Args *[]string) {
		if err := kindApp.Run(cmd.NewLogger(), cmd.StandardIOStreams(), *Args); err != nil {
			os.Exit(1)
		}
		wg.Done()
	}

	for i := 0; i < len(karray.Cluster); i++ {
		for j := 0; j < len(karray.Cluster[i].ImagesToLoad); j++ {
			aux := []string{"load"}
			aux = append(aux, "docker-image", "--name", karray.Cluster[i].Name, karray.Cluster[i].ImagesToLoad[j].Image)
			aux = append(aux, karray.Cluster[i].ImagesToLoad[j].Image)
			fmt.Println("---> kind ", aux)
			wg.Add(1)
			go KindExec(&aux)
		}

	}
}

func clab(mytopofile string) {

	cmd := exec.Command("/usr/bin/containerlab", "deploy", "--topo", mytopofile)
	cmd.Stdout = &stdOut{}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func cklinks(mylinks []Link_Data) {

	for i := 0; i < len(mylinks); i++ {

		cmd := exec.Command("/usr/bin/containerlab", "tools", "veth", "create", "-a", mylinks[i].K8sNode, "-b", mylinks[i].ClabNode)
		fmt.Printf("---> /usr/bin/containerlab tools veth create -a %s -b %s\n", mylinks[i].K8sNode, mylinks[i].ClabNode)
		cmd.Stdout = &stdOut{}
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func kindsetInf(mylinks []Link_Data) {

	for i := 0; i < len(mylinks); i++ {
		re1 := regexp.MustCompile(`e\d+-\d+`)
		myif := re1.FindString(mylinks[i].K8sNode)
		re2 := regexp.MustCompile(`:e\d+-\d+`)
		mynode := re2.ReplaceAllString(mylinks[i].K8sNode, "")
		//KNetScript(mylinks[i].K8sIpv4, myif, mylinks[i].K8sIpv4Gw, mynode)
		Bond0Create(bondif, myif, mynode, mylinks[i].K8sIpv4, mylinks[i].K8sIpv4Gw)
		for j := 0; j < len(mylinks[i].IpvlanMaster); j++ {
			KNetVlanMaster(mylinks[i].IpvlanMaster[j].Vlan, bondif, myif, mynode)

		}
	}

}

func execApps(myres []Resources_Data, kubeconf string, myns string) {
	for i := 0; i < len(myres); i++ {
		ExecKubeApps(myres[i], kubeconf, myns)
	}
	wg.Done()
}
