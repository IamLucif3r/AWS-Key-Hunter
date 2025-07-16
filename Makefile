.PHONY: all build docker clean

BINARY_NAME=aws-key-scanner
DOCKER_IMAGE_NAME=aws-key-hunter:latest

all: build docker

build:
	@echo "🔨 Building Go binary..."
	go build -o $(BINARY_NAME) ./cmd/awsKeyhunter.go

docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME)
