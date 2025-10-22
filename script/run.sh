#!/bin/bash

echo "welcome to pasunes"

pwd="../cmd/pasunes_scan/main.go"

target="172.17.0.1"


go run $pwd -target $target -v
