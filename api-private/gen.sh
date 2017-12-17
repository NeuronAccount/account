#!/usr/bin/env bash

rm -rf ./gen/
mkdir gen

swagger generate server -T ~/work/neuron/src/github.com/NeuronFramework/restful/go_template/ -f swagger.json -t ./gen/