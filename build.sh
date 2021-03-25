#! /bin/bash

export VERSION=$(git describe --tags --long --abbrev=12 --match "v[0-9].*" --always)
go build -o ./kubectl-opslevel -ldflags="-X 'github.com/opslevel/kubectl-opslevel/cmd.version=${VERSION}'"
