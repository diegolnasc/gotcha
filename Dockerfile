ARG BINARY_FOLDER=bin
ARG BINARY_NAME=gotcha
ARG GOOS=linux
ARG GOARCH=amd64

FROM golang:1.18 AS builder

ARG BINARY_FOLDER
ARG BINARY_NAME
ARG GOOS
ARG GOARCH

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o $BINARY_FOLDER/$BINARY_NAME .

FROM gcr.io/distroless/base-debian10

ARG BINARY_FOLDER
ARG BINARY_NAME

WORKDIR /

COPY --from=builder app/build/ build/
COPY --from=builder app/$BINARY_FOLDER/$BINARY_NAME gotcha

EXPOSE 3000
ENTRYPOINT [ "/gotcha", "server"]
