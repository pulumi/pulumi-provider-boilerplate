#!/bin/bash
npm install -g yarn
pip install setuptools

curl --silent -L https://github.com/pulumi/pulumictl/releases/download/v0.0.47/pulumictl-v0.0.47-linux-amd64.tar.gz | sudo tar -xz -C /usr/local/bin pulumictl
sudo chmod +x /usr/local/bin/pulumictl