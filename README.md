# Scripting with pyGNMI to troubleshoot Datacenters

Welcome to this tutorial where we will guide you through the creation of a network lab using containerlab. Specifically, we will be focusing on the Spine/Leaf + Border-Leaf Datacenter Architecture, using technologies such as Layer 2/3 EVPN-VXLAN and eBGP for Underlay.

We will cover everything you need to know about creating a virtual network environment that emulates a real-world network, complete with routing protocols and traffic flow. You will learn how to install the necessary requirements for pyGNMI scripting, a powerful tool that allows you to operate and troubleshoot EVPN-VXLAN design.

Other important topic added to this tutorial is GNMIc, in conjunction with Go templates. Those allows us to orchestrate configurations in network devices by translating instructions from a YAML file. This process involves reading the variable values from the YAML file and applying the Go template to construct a GNMI request. It's a powerful tool that offers flexibility and a useful alternative to mainstream automation frameworks like Ansible, particularly when managing complex configurations or operating on multiple devices simultaneously.

By the end of this tutorial, you will be able to script in pyGNMI to perform various tasks in your virtual network environment. This will give you a deeper understanding of how networks operate and enable you to troubleshoot and optimize network performance. Whether you are an experienced network engineer or just starting, this tutorial will provide valuable insights into network design and operation. Let's get started!

## Installing containerlab

This script installs and starts Docker, a containerization platform, on a Linux machine using the dnf package manager. Then, it installs containerlab, a tool used for creating and managing container-based network labs. The command specified installs containerlab version 0.25.1. (it's for Fedora33)

```
# Install docker
sudo dnf -y install docker
sudo systemctl start docker
sudo systemctl enable docker

# Install containerlab
bash -c "$(curl -sL https://get.containerlab.dev)" -- -v 0.25.1
```

## Installing pyGNMI and requirements

Necessary Python packages for your script including `pygnmi`, `prettytable`, `yaml` and others.

### Pre-requisites

You will need Python 3.6 or newer installed on your machine. You can download Python from the [official website](https://www.python.org/downloads/).


#### Step 1: Install pip

Pip is the package installer for Python. If you don't have pip installed, you can download and install it from [here](https://pip.pypa.io/en/stable/installing/).

#### Step 2: Install necessary Python packages

Run the following commands to install the necessary packages:

```bash
pip install pygnmi
pip install prettytable
pip install PyYAML
pip install tabulate
```

Note: The operator and time packages are part of the Python Standard Library and should be available by default, so you don't need to install them separately.

#### Step 3: Import packages
Once you've installed the packages, you can import them in your Python script as follows:

```python
from pygnmi.client import gNMIclient
from prettytable import PrettyTable
import yaml
import time
import logging
from operator import itemgetter
from itertools import groupby
from tabulate import tabulate

```

## Installing GNMIc

It should be as easy as:
```bash
bash -c "$(curl -sL https://get-gnmic.kmrd.dev)"
```

You can use the following link for more details: https://gnmic.kmrd.dev/install/


## Configuring network elements via GNMIc
Other important topic added to this tutorial is GNMIc, in conjunction with Go templates. Those allows us to orchestrate configurations in network devices by translating instructions from a YAML file. This process involves reading the variable values from the YAML file and applying the Go template to construct a GNMI request. It's a powerful tool that offers flexibility and a useful alternative to mainstream automation frameworks like Ansible, particularly when managing complex configurations or operating on multiple devices simultaneously.

We'll use go templates with the following format (in this case to create or replace network instances and other elements):
```go
replaces:
{{ $target := index .Vars .TargetName }}
{{- range $netinstances := index $target "network-instances" }}
  - path: "/network-instance[name={{ index $netinstances "name" }}]"
    encoding: "json_ietf"
    value: 
      admin-state: {{ index $netinstances "admin-state" | default "disable" }}
      type: {{ index $netinstances "type" | default "mac-vrf" }}
      description: {{ index $netinstances "description" | default "whatever" }}
      vxlan-interface:
        - name: vxlan2.{{ index $netinstances "vni" }}
      protocols:
        bgp-evpn:
          bgp-instance:
            - id: 2
              admin-state: enable
              vxlan-interface: vxlan2.{{ index $netinstances "vni" }}
              evi: {{ index $netinstances "evi" }}
        bgp-vpn:
          bgp-instance:
            - id: 2
              route-target:
                export-rt: target:65123:{{ index $netinstances "evi" }}
                import-rt: target:65123:{{ index $netinstances "evi" }}                           
{{- end }}
```

and then we can use YML file defining the vars with a format like the following:
```yaml
clab-dc-k8s-LEAF-DC-1:
  network-instances:
    - name: l2evpn1001
      admin-state: enable
      type: mac-vrf
      evi: 1001
      vni: 1001
      vxtype: bridged
      anycast-gw: 10.0.1.1/24
    - name: l2evpn1002
      admin-state: enable
      type: mac-vrf
      evi: 1002
      vni: 1002
      vxtype: bridged
      anycast-gw: 10.0.2.1/24
```


To see the actual files we'll use on this tutorial check at https://github.com/cloud-native-everything/pygnmi-srl-nanog88/tree/main/gnmic

## Checking EVPN settings using pyGNMI

This Python application connects to a list of specified routers using the gNMI protocol and retrieves both BGP EVPN and BGP VPN information. The information is then formatted into a PrettyTable for easy viewing.

The application primarily consists of the SrlDevice class, which represents a router. This class is initialized with the router's basic information and uses the gNMI client to fetch BGP EVPN and BGP VPN information.

The application then creates a list of these SrlDevice instances based on a YAML configuration file ('datacenter-nodes.yml'). It generates two tables: one sorted by router name and another sorted by network instance.

```bash
[root@rbc-r2-hpe4 py-scripts]# python3 display_evpn_per_router.py datacenter-nodes.yml
Table: Sorted by Router
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
| Router                | Network instance | ID | EVPN Admin state | VXLAN interface | EVI  | ECMP | Oper state | RD           | import-rt         | export-rt         |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
| clab-dc-k8s-LEAF-DC-1 | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.1:4    | target:65123:4    | target:65123:4    |
| clab-dc-k8s-LEAF-DC-1 | kube_macvrf      | 1  | enable           | vxlan1.1        | 1    | 1    | up         | 1.1.1.1:1    | target:65123:1    | target:65123:1    |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1001       | 2  | enable           | vxlan2.1001     | 1001 | 1    | no state   | 1.1.1.1:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1002       | 2  | enable           | vxlan2.1002     | 1002 | 1    | no state   | 1.1.1.1:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1003       | 2  | enable           | vxlan2.1003     | 1003 | 1    | no state   | 1.1.1.1:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1004       | 2  | enable           | vxlan2.1004     | 1004 | 1    | no state   | 1.1.1.1:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1005       | 2  | enable           | vxlan2.1005     | 1005 | 1    | no state   | 1.1.1.1:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-1 | l2evpn1006       | 2  | enable           | vxlan2.1006     | 1006 | 1    | no state   | 1.1.1.1:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-1 | l3evpn           | 1  | enable           | vxlan1.2        | 2    | 4    | up         | 1.1.1.1:2    | target:65123:2    | target:65123:2    |
| clab-dc-k8s-LEAF-DC-2 | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.2:4    | target:65123:4    | target:65123:4    |
| clab-dc-k8s-LEAF-DC-2 | kube_macvrf      | 1  | enable           | vxlan1.1        | 1    | 1    | up         | 1.1.1.2:1    | target:65123:1    | target:65123:1    |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1001       | 2  | enable           | vxlan2.1001     | 1001 | 1    | up         | 1.1.1.2:1001 | target:65123:1001 | target:65123:1001 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1002       | 2  | enable           | vxlan2.1002     | 1002 | 1    | up         | 1.1.1.2:1002 | target:65123:1002 | target:65123:1002 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1003       | 2  | enable           | vxlan2.1003     | 1003 | 1    | no state   | 1.1.1.2:1003 | target:65123:1003 | target:65123:1003 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1004       | 2  | enable           | vxlan2.1004     | 1004 | 1    | no state   | 1.1.1.2:1004 | target:65123:1004 | target:65123:1004 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1005       | 2  | enable           | vxlan2.1005     | 1005 | 1    | no state   | 1.1.1.2:1005 | target:65123:1005 | target:65123:1005 |
| clab-dc-k8s-LEAF-DC-2 | l2evpn1006       | 2  | enable           | vxlan2.1006     | 1006 | 1    | no state   | 1.1.1.2:1006 | target:65123:1006 | target:65123:1006 |
| clab-dc-k8s-LEAF-DC-2 | l3evpn           | 1  | enable           | vxlan1.2        | 2    | 4    | up         | 1.1.1.2:2    | target:65123:2    | target:65123:2    |
| clab-dc-k8s-BORDER-DC | kube-ipvrf       | 1  | enable           | vxlan1.4        | 4    | 4    | up         | 1.1.1.10:4   | target:65123:4    | target:65123:4    |
+-----------------------+------------------+----+------------------+-----------------+------+------+------------+--------------+-------------------+-------------------+
Total time: 1.52 seconds
```

To see the actual files we'll use on this tutorial check at https://github.com/cloud-native-everything/pygnmi-srl-nanog88/tree/main/py-scripts