#!/bin/sh
set -e

files=$(find -type f -name "*.proto")

for proto in $files; do
    protoc $proto --proto_path=$(dirname $proto) --go_out=. --go-grpc_out=.
done
