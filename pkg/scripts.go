package srklab

import (
	"fmt"
	"os/exec"
	"strings"
)

func KNetScript(myipv4 string, myif string, mygw string, mynode string) {
	cmd := exec.Command("docker", "exec", "-i", mynode, "bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`#!/bin/bash	
	ip link set %s down
	ip link set %s up mtu 9500
	ip addr add %s dev %s
	ip link set %s up
	ip route del default
	ip route add default via %s
	`, myif, myif, myipv4, myif, myif, mygw))
	fmt.Printf("INFO: Executing bash scripts to change settings for node %s and interface %s\n", mynode, myif)
	output, err := cmd.CombinedOutput()
	fmt.Println("-->", string(output))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func KNetVlanMaster(vlan string, bondif string, myif string, mynode string) {
	cmd := exec.Command("docker", "exec", "-i", mynode, "bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`#!/bin/bash	
	ip link add link %s name VLAN-%s type vlan id %s
    ip link set dev VLAN-%s up
	`, bondif, vlan, vlan, vlan))
	fmt.Printf("INFO: Executing bash scripts to create vlan %s at interface %s\n", vlan, myif)
	output, err := cmd.CombinedOutput()
	fmt.Println("-->", string(output))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func Bond0Create(bondif string, myif string, mynode string, myipv4 string, mygw string) {
	cmd := exec.Command("docker", "exec", "-i", mynode, "bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`#!/bin/bash	
	ip link add %s type bond
	ip link set %s down
	ip link set %s master %s
	ip link set %s up mtu 9500
	ip addr add %s dev %s
	ip route replace default via %s	
	`, bondif, myif, myif, bondif, bondif, myipv4, bondif, mygw))
	fmt.Printf("INFO: Executing bash scripts to create bond interface %s, with IP %s, Gateway %s, at interface %s\n", bondif, myipv4, mygw, myif)
	output, err := cmd.CombinedOutput()
	fmt.Println("-->", string(output))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}
