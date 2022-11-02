#!/bin/bash
protoc -I=. -I=../.. --go_opt=paths=source_relative --go_out=. dataplane.proto
protoc -I=. -I=../.. --go-grpc_opt=paths=source_relative --go-grpc_out=. dataplane.proto
