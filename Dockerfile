FROM golang:alpine AS builder

ENV BINARY_FOLDER=bin
ENV BINARY_NAME=gotcha
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR $GOPATH/gotcha
COPY . $GOPATH/gotcha/
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN go fmt ./cmd/ ./pkg/config/ ./pkg/github/
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o ${BINARY_FOLDER}/${BINARY_NAME} ./cmd/main.go
EXPOSE 3000
ENTRYPOINT "${BINARY_FOLDER}/${BINARY_NAME}"