FROM golang:1 as builder

LABEL writer="zerrozhao@gmail.com"

ENV GOPROXY https://goproxy.io

WORKDIR /src/tradingdb2

COPY ./go.* /src/tradingdb2/

RUN go mod download

COPY . /src/tradingdb2

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tradingdb2 ./node \
    && mkdir /app \
    && mkdir /app/tradingdb2 \
    && mkdir /app/tradingdb2/cfg \
    && mkdir /app/tradingdb2/logs \
    && mkdir /app/tradingdb2/data \
    && cp ./tradingdb2 /app/tradingdb2/ \
    && cp ./VERSION /app/tradingdb2/ \
    && cp -r ./cfg /app/tradingdb2/ \
    && cp ./cfg/config.yaml.default /app/tradingdb2/cfg/config.yaml

FROM alpine
RUN apk upgrade && apk add --no-cache ca-certificates
WORKDIR /app/tradingdb2
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /app/tradingdb2 /app/tradingdb2
CMD ["./tradingdb2"]