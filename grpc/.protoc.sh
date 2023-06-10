#!/bin/sh

protoc \
  --proto_path .protos/v1 \
  --go_opt paths=source_relative \
  --go_out v1 \
  --go-grpc_opt paths=source_relative \
  --go-grpc_out v1 \
  system.proto
