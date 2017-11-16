#!/usr/bin/env bash

rm -rf ./gen/
mkdir gen

swagger generate server -f swagger.json -t ./gen/
swagger generate client -f swagger.json -t ./gen/