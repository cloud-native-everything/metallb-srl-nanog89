package srklab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"text/template"
)

type Hostname struct {
	Hostname string
	ipaddr   string
}

func Metrics(configFile string) {

	// Read the template file into a string
	templateFile, err := ioutil.ReadFile("./prometheus/prometheus.tmpl")
	if err != nil {
		panic(err)
	}
	hosts := []Hostname{
		{Hostname: "edge1-control-plane", ipaddr: "172.18.45.45"},
		{Hostname: "netorch-control-plane", ipaddr: "172.18.46.46"},
	}

	// Define the template
	tmpl, err := template.New("person").Parse(string(templateFile))
	if err != nil {
		panic(err)
	}

	// Render the template with data
	var output bytes.Buffer
	err = tmpl.Execute(&output, hosts)
	if err != nil {
		panic(err)
	}

	// Write the output to a file
	err = ioutil.WriteFile("./prometheus/prometheus.cfg", output.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Output written to ./prometheus/prometheus.cfg")

	kindVars := getVars(configFile)
	host_list := &KHostList
	for i := 0; i < len(kindVars.Cluster); i++ {
		MyNodes(kindVars.Cluster[i].Name, host_list)
	}
	fmt.Println(*host_list)
}
func MyNodes(mycluster string, host_list *[]Kind_Hosts) {

	cmd := exec.Command("/root/go/bin/kind", "get", "nodes", "--name", mycluster)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	files := strings.Split(string(out), "\n")
	filesSlice := make([]string, 0, len(files))
	for _, file := range files {
		if file != "" {
			*host_list = append(*host_list, getIP(file))
			filesSlice = append(filesSlice, file)
		}
	}

}

func getIP(container string) Kind_Hosts {
	cmd := exec.Command("docker", "inspect", "-f", "'{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}'", container)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	ipadd := strings.TrimSuffix(strings.ReplaceAll(string(out), "'", ""), "\n")
	fmt.Printf("-------> %s -- %s\n", ipadd, container)
	return Kind_Hosts{Hostname: container, Ipadd: ipadd}

}
