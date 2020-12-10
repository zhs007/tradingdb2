rm -rf grpc/tradingdb2
rm -rf grpc/jrj
rm -rf grpc/bitmet
rm -rf grpc/simtrading
rm -rf unittestdata/tradingdb2
rm -rf unittestdata/jrj
rm -rf unittestdata/bitmet
rm -rf unittestdata/simtrading
rm -rf utils/*.log
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out