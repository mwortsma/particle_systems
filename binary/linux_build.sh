#!/bin/bash
env GOOS=linux GOARCH=arm go build -v github.com/mwortsma/particle_systems/dtlb/dtlb
env GOOS=linux GOARCH=arm go build -v github.com/mwortsma/particle_systems/dtcp/dtcp
