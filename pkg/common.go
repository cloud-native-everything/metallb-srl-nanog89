package srklab

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

type stdOut struct{}

var wg = sync.WaitGroup{}

var KHostList []Kind_Hosts

const bondif = "bond0"

type Kind_Data struct {
	Cluster      []Cluster_Data `yaml:"clusters"`
	Links        []Link_Data    `yaml:"links"`
	Network      string         `yaml:"network"`
	Prefix       string         `yaml:"prefix"`
	ClabTopology string         `yaml:"clabTopology"`
}

type Cluster_Data struct {
	//Info of every K8s Kind Cluster
	Name         string              `yaml:"name"`
	Config       string              `yaml:"config"`
	Kubeconfig   string              `yaml:"kubeconfig"`
	Image        string              `yaml:"image"`
	ImagesToLoad []ImagesToLoad_Data `yaml:"imagesToLoad"`
	Resources    []Resources_Data    `yaml:"resources"`
}

type Link_Data struct {
	K8sNode      string      `yaml:"k8sNode"`
	ClabNode     string      `yaml:"clabNode"`
	K8sIpv4      string      `yaml:"k8sIpv4"`
	K8sIpv4Gw    string      `yaml:"k8sIpv4Gw"`
	IpvlanMaster []Vlan_Data `yaml:"ipvlanMaster"`
}

type Vlan_Data struct {
	Vlan string `yaml:"vlan"`
}

type ImagesToLoad_Data struct {
	Image string `yaml:"image"`
}

type Resources_Data struct {
	App string `yaml:"app"`
}

type Kind_Hosts struct {
	Hostname string
	Ipadd    string
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getVars(configFile string) (kdata *Kind_Data) {
	gtt_config, err := ioutil.ReadFile(configFile)
	checkerr(err)
	err = yaml.Unmarshal(gtt_config, &kdata)
	return kdata

}

func (s *stdOut) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}
