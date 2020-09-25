export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=proto/ --go_out=plugins=grpc:tradingdb2pb --go_opt=paths=source_relative proto/*.proto