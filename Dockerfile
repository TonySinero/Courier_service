FROM golang:alpine AS builder

COPY . /github.com/Baraulia/COURIER_SERVICE/
WORKDIR /github.com/Baraulia/COURIER_SERVICE/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/Baraulia/COURIER_SERVICE/.bin/service .
COPY --from=0 /github.com/Baraulia/COURIER_SERVICE/configs configs/

EXPOSE 81 8091

CMD ["./service"]

