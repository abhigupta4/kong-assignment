GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=kong
BUILD_DIR=.

build:
	$(GOBUILD) -o $(BINARY_NAME) $(BUILD_DIR)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	$(GOCMD) fmt ./...

run-kafka-producer: build
	APP_TYPE=kafka-producer ./$(BINARY_NAME)

run-kafka-consumer: build
	APP_TYPE=kafka-consumer ./$(BINARY_NAME)

all: clean fmt build