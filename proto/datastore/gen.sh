#!/bin/bash
protoc --proto_path=. --proto_path=${GOPATH}/src --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative datastore.proto
