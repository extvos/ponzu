#!/bin/bash

# Set up test environment

set -ex 

# Install Ponzu CMS
go get -u github.com/extvos/ponzu/...


# test install
ponzu


# create a project and generate code
ponzu new github.com/extvos/ci/test-project

cd /go/src/github.com/extvos/ci/test-project

ponzu gen content person name:string hashed_secret:string
ponzu gen content message from:@person,hashed_secret to:@person,hashed_secret


# build and run dev http/2 server with TLS
ponzu build

