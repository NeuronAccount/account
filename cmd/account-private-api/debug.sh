#!/usr/bin/env bash

mkdir _output
rm -rf ./_output/debug
mkdir ./_output/debug

go build -o ./_output/debug/main .

ENV=DEV \
DEBUG=true \
PORT=8083 \
DB="root:123456@tcp(127.0.0.1:3307)" \
./_output/debug/main