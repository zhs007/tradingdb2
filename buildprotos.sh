export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=protos/ --go_out=./tradingpb --go_opt=paths=source_relative protos/*.proto
protoc --proto_path=protos/ --go-grpc_out=./tradingpb --go-grpc_opt=paths=source_relative protos/*.proto