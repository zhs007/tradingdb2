rm -rf unittestdata/tradingdb2
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out