FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .

#增加缺失的包，移除没用的包
RUN go mod tidy
RUN go build -o main .
# 下载最新的数据库
# RUN wget https://github.do/https://raw.githubusercontent.com/bqf9979/ip2region/master/data/ip2region.db


FROM scratch
# FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /app/ip2region.db /app
EXPOSE 9090
ENTRYPOINT  ["/app/main"]