#!/bin/bash
#
# ./setup_protogen.sh
#
# Install protobuf and protogen related go libs
function setup_go_protogen_execs() {
  export GOBIN=$(go env GOBIN) 

  go mod download google.golang.org/protobuf
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
}


setup_go_protogen_execs
