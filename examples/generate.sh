#!/bin/bash
#
# Generate examples
#
#
mkdir -p httponly
mkdir -p withgrpc
mkdir -p openapiv3

### generate httponly
cd httonly;go mod init testhttpsetup;bananas init --name=testhttpsetup --oa-version=2;bananas gen:docs --oa-version=2;cd ..

### generate grpc server
cd withgrpc;go mod init testgrpcsetup;bananas init --name=testgrpcsetup --grpc --oa-version=2;cd ..


### generate openapiv3
cd openapiv3;go mod init openapivvv;bananas init --name=openapivvv;bananas gen:docs;cd ..
