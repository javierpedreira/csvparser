#!/usr/bin/env bash

set -ex

docker run --rm -e GOOS=darwin -v "${PWD}":/parser -w /parser golang:1.14 go build -v
