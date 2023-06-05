package DockerNet

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func DockerNetworkCreate(mynetwork string, myprefix string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Create an IPAM configuration with a subnet
	ipamConfig := &network.IPAMConfig{
		Subnet: myprefix,
	}

	networkResponse, err := cli.NetworkCreate(context.Background(), mynetwork, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "bridge",
		EnableIPv6:     false,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{*ipamConfig},
		},
		Internal: false,
	})
	if err != nil {
		panic(err)
	}
	if networkResponse.ID != "network_id" {
		fmt.Printf("expected networkResponse.ID to be 'network_id', got %s\n", networkResponse.ID)
	}
	if networkResponse.Warning != "warning" {
		fmt.Printf("expected networkResponse.Warning to be 'warning', got %s\n", networkResponse.Warning)
	}
}

func DockerNetworkDelete( mynetwork string) {
	var networkID string 
        
	cli, err := client.NewEnvClient()
	if err != nil {
    		panic(err)
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
	    panic(err)
	}

	for _, n := range networks {
	    if n.Name == mynetwork {
      			networkID = n.ID
        	// Do something with the network ID
  		}
	}


        err = cli.NetworkRemove(context.Background(), networkID)
	if err != nil {
        	panic(err)
	}

}
