#!/bin/bash
set -eu

dockerd-entrypoint.sh

go test -v ./...