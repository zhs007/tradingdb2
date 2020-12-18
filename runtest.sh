rm -rf grpc/tradingdb2
rm -rf grpc/jrj
rm -rf grpc/bitmex
rm -rf grpc/simtrading
rm -rf unittestdata/tradingdb2
rm -rf unittestdata/jrj
rm -rf unittestdata/bitmex
rm -rf unittestdata/simtrading
rm -rf utils/*.log

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

rm -rf grpc/tradingdb2
rm -rf grpc/jrj
rm -rf grpc/bitmex
rm -rf grpc/simtrading
rm -rf unittestdata/tradingdb2
rm -rf unittestdata/jrj
rm -rf unittestdata/bitmex
rm -rf unittestdata/simtrading
rm -rf utils/*.log