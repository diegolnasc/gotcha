BINARY_NAME=gotcha
BINARY_FOLDER=bin
BINARY_NAME=gotcha
GOOS=linux
GOARCH=amd64

build: clean
	mkdir ${BINARY_FOLDER} && \
	go mod tidy && \
	go mod download && \
	export GO111MODULE=on && \
	env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o ${BINARY_FOLDER}/${BINARY_NAME} .

docker_build:
	docker build -t ${BINARY_NAME} .

docker_run:
	docker run -p 3000:3000 --name ${BINARY_NAME} ${BINARY_NAME}

docker_all: docker_build docker_run

run: 
	./${BINARY_FOLDER}/${BINARY_NAME} server

clean:
	go clean	
	rm -rf ${BINARY_FOLDER}/