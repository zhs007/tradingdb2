rm -rf unittestdata/tradingdb2
rm -rf utils/*.log
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out