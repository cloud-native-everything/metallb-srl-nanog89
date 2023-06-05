#!/bin/bash

# Install docker
if command -v dnf >/dev/null 2>&1; then
  sudo dnf -y install docker bridge-utils
elif command -v apt-get >/dev/null 2>&1; then
  sudo apt-get update
  sudo apt-get -y install docker.io
  sudo apt-get -y install bridge-utils
else
  echo "Unsupported package manager, please install docker manually."
  exit 1
fi

sudo systemctl start docker
sudo systemctl enable docker

# Install containerlab
bash -c "$(curl -sL https://get.containerlab.dev)" -- -v 0.25.1

# Install go
if command -v dnf >/dev/null 2>&1; then
  sudo dnf update -y
elif command -v apt-get >/dev/null 2>&1; then
  sudo apt-get update
fi

curl -LO https://golang.org/dl/go1.17.8.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.17.8.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
rm -f go1.17.8.linux-amd64.tar.gz
go version

