#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -v github.com/mwortsma/particle_systems/dtlb/dtlb
env GOOS=linux GOARCH=amd64 go build -v github.com/mwortsma/particle_systems/dtcp/dtcp
