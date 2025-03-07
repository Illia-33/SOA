protoc $PWD/proto/*.proto --proto_path=$PWD/proto --go_out=./src --go-grpc_out=./src
go build $PWD/src
