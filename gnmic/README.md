# GNMIc and Go Templates Tutorial

GNMIc is an open-source GNMI (gRPC Network Management Interface) client written in Go language. It is a versatile tool that can be used to manage network devices using the YANG models and GNMI standard defined by the OpenConfig working group. GNMIc can perform CRUD (Create, Read, Update, Delete) operations and also subscribe to telemetry streams. This tutorial will guide you through using GNMIc command line interface (CLI) and how to leverage Go templates for complex tasks.

## Example CLI usage

In the following example, GNMIc is used to perform a simple read operation on a network device to retrieve the hostname:

```bash
gnmic -a 172.18.100.122:57400 -u admin -p admin --skip-verify get -e json_ietf --path /system/name/host-name
```

In this command:

* `-a 172.18.100.122:57400` specifies the address (IP and port) of the target device.
* `-u admin -p admin` provide the username and password for authentication.
* `--skip-verify` is used to skip server certificate validation, useful in development environments.
* `get` is the action to perform (a read operation in this case).
* `-e json_ietf` is the encoding format. Here, it is set to json_ietf, indicating that the server should return data in IETF JSON format.
* `--path /system/name/host-name` is the path of the data in the YANG model that we want to retrieve.

This command would return the hostname of the network device at the specified address in JSON format.

## Go Templates for creating network instances and interfaces
Go Templates can be used along with GNMIc to simplify the process of creating or modifying multiple network instances and interfaces. Here is an example:

```bash
gnmic -a clab-dc-k8s-LEAF-DC-1,clab-dc-k8s-LEAF-DC-2 -u admin -p admin --skip-verify set --request-vars l2-evpn-domains-create-vars.yml --request-file l2-evpn-domains-create.gotmpl
```
In this command:

* `-a clab-dc-k8s-LEAF-DC-1,clab-dc-k8s-LEAF-DC-2` specifies the addresses of the target devices. Here, we are targeting two devices.
* `set` is the action to perform (a create or update operation in this case).
* `--request-vars l2-evpn-domains-create-vars.yml` is a file that contains variables that will be used in the Go template.
* `--request-file l2-evpn-domains-create.gotmpl` is the Go template file that will be used to construct the request.

This command would read the variable values from the specified YAML file and apply the Go template to construct a GNMI set request. This is sent to each of the network devices specified in 
the command, effectively creating or updating the L2 EVPN domains on these devices as per the data in the template.

Go Templates provide a powerful and flexible way to manage network configuration, especially when dealing with complex configurations or multiple devices.

more details at: https://gnmic.kmrd.dev/