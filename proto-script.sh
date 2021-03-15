protoc -I . \
  --go_out . \
  --go_opt paths=source_relative \
  --go-grpc_out . --go-grpc_opt paths=source_relative \
  proto/todopb/todo_server.proto

protoc -I . \
    --grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
  proto/todopb/todo_server.proto
