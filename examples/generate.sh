#!/bin/bash
#
# Generate examples
#
#
mkdir -p httponly
mkdir -p withgrpc

### generate httponly
cd httonly;go mod init testhttpsetup;bananas init --name=testhttpsetup;bananas gen:docs;cd ..

### generate grpc server
cd withgrpc;go mod init testgrpcsetup;bananas init --name=testgrpcsetup --grpc;bananas gen:docs;cd ..
