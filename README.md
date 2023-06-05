# Kubernetes BGP True Load Balancer for Datacenters

Are you struggling with load balancing in your on-premises Kubernetes cluster? Do you wish to have the same level of automation and experience as the Public Cloud? Look no further! In this presentation, we will guide you through defining your own on-premises Kubernetes LoadBalancer service using BGP through the Datacenter Fabric and bringing true load balancing across the leaf switches with ECMP.

We will demonstrate how to set up a demo from scratch using open-source tools like Containerlab, MetalLB, and Kubernetes Kind. MetalLB is one of the most widely used open-source load balancer projects in enterprises and is suitable for telco use cases like IoT or 5G edge designs. Kubernetes Kind is a tool for running local Kubernetes clusters using Docker container “nodes” and can be used for local development or CI.
This presentation is intended for audiences with any level of skills.

We welcome any collaboration on this project.

Don't miss this opportunity to learn how to bring true load balancing to your on-premises Kubernetes cluster.

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

