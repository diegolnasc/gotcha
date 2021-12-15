BINARY_NAME=gotcha

build: clean
	mkdir bin
	go mod tidy && \
	go mod download && \
	export GO111MODULE=on && \
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME} ./cmd/main.go

docker_build:
	docker build -t ${BINARY_NAME} .

docker_run:
	docker run -p 3000:3000 ${BINARY_NAME} --name ${BINARY_NAME}

docker_all: docker_build docker_run

run: 
	./bin/${BINARY_NAME}

clean:
	go clean
	rm -rf bin/