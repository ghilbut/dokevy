#!/bin/sh

protoc \
  --proto_path ./grpc/.protos/v1 \
  --go_opt paths=source_relative \
  --go_out ./grpc/v1 \
  --go-grpc_opt paths=source_relative \
  --go-grpc_out ./grpc/v1 \
  --plugin protoc-gen-ts_proto=./next.js/node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_opt outputServices=grpc-js,env=node,esModuleInterop=true \
  --ts_proto_out ./next.js/grpc/v1 \
  system.proto
