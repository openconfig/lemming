protoc --go_opt=paths=source_relative --go_out=. dataplane/*.proto
protoc --go-grpc_opt=paths=source_relative --go-grpc_out=. dataplane/dataplane.proto