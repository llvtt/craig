# Generic Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=craig-server
BINARY_UNIX=$(BINARY_NAME)_unix

# Craig parameters
ENV_FILE='./.env'
CONFIG_FILE='./dev.config.json'
include .env
export

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v main/main.go
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN) ./main
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
#		if [ -f $(ENV_FILE) ]; then
#		fi
		./$(BINARY_NAME) --config-file=$(CONFIG_FILE)
deps:
		go mod download
		go mod verify


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
        #TODO
		#docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v