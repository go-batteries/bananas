
#!/bin/bash
#
# ./setup_apidocs.sh
#
# Install protobuf and protogen related go libs
function setup_go_protogen_execs() {
  export GOBIN=$(go env GOBIN) 

  if [[ -z "$GOBIN" ]]; then
    printf "Fatal!! \$GOBIN is not set in current environment.\nSet \$GOBIN to appropriate directory and add to \$PATH\n"
  fi

  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
  go mod download google.golang.org/protobuf
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
  go install github.com/google/gnostic/cmd/protoc-gen-openapi
  go install github.com/go-swagger/go-swagger/cmd/swagger@latest
}


setup_go_protogen_execs
