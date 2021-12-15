FROM golang:alpine AS builder
COPY . $GOPATH/gotcha/
WORKDIR $GOPATH/gotcha
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/gotcha ./cmd/main.go
EXPOSE 3000
ENTRYPOINT ["./bin/gotcha"] 