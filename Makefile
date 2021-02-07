# Generic Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
LINUXFLAGS=GOOS=linux GOARCH=amd64
TERRAFORM_DIR=./terraform
SLACK_EVENTS_BINARY_NAME=slack-events
CLOUDWATCH_EVENTS_BINARY_NAME=cloudwatch-events

all: test build
build-dev: clean deps
		$(GOBUILD) -o $(SLACK_EVENTS_BINARY_NAME)-dev -v cmd/slack-events/main.go
		$(GOBUILD) -o $(CLOUDWATCH_EVENTS_BINARY_NAME)-dev -v cmd/cloudwatch-events/main.go

build: clean deps
		$(LINUXFLAGS) $(GOBUILD) -o $(SLACK_EVENTS_BINARY_NAME) -v cmd/slack-events/main.go
		$(LINUXFLAGS) $(GOBUILD) -o $(CLOUDWATCH_EVENTS_BINARY_NAME) -v cmd/cloudwatch-events/main.go

test:
		$(GOTEST) -v ./...

clean:
		$(GOCLEAN) ./cmd
		rm -f $(BINARY_NAME)

clean-lambda:
		$(GOCLEAN) ./main
		rm -f $(LAMBDA_DIR)/$(BINARY_NAME)

run:
		./$(BINARY_NAME)

deps:
		go mod download
		go mod verify

docker-build: build
		docker build -t "craig" .

docker--set-tag-name:
		$(eval tag_name=$(shell docker images --no-trunc --quiet craig | cut -d: -f2))

docker-push: docker-build docker--set-tag-name
		docker tag "craig:latest" "${ECR_HOSTNAME}/craig:${tag_name}"
		docker run --rm -it -e AWS_DEFAULT_REGION -e  AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY amazon/aws-cli ecr get-login-password | docker login --username AWS --password-stdin ${ECR_HOSTNAME}
		docker push "${ECR_HOSTNAME}/craig:${tag_name}"

deploy-plan: docker--set-tag-name
		cd $(TERRAFORM_DIR) && terraform init && terraform plan -var="tag_name=${tag_name}"

deploy: docker--set-tag-name
		cd $(TERRAFORM_DIR) && terraform init && terraform apply -auto-approve -var="tag_name=${tag_name}"
