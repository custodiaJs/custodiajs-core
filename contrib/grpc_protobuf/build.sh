export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=../../localgrpcproto --go_opt=paths=source_relative \
       --go-grpc_out=../../localgrpcproto --go-grpc_opt=paths=source_relative \
       local_rpc.proto