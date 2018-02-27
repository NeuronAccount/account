#!/usr/bin/env bash

mkdir _output
rm -rf ./_output/release
mkdir ./_output/release

GOOS=linux GOARCH=amd64 go build -o ./_output/release/neuron-agent .

docker build -t neuron-account-private-api:v1.0.0 .