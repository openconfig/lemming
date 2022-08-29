#!/bin/bash
cd ../..
protoc --go_opt=paths=source_relative --go_out=. proto/forwarding/*.proto
protoc --go-grpc_opt=paths=source_relative --go-grpc_out=. proto/forwarding/*.proto