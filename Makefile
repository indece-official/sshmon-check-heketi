# Project config
PROJECT_NAME=sshmon_check_heketi

# Go parameters
GOCMD=go
GOPATH=$(shell $(GOCMD) env GOPATH))
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DIR_SOURCE=./src
DIR_DIST=./dist
BINARY_NAME=$(DIR_DIST)/bin/$(PROJECT_NAME)
BUILD_DATE=$(shell date +%Y%m%d.%H%M%S)
BUILD_VERSION ?= $(shell git rev-parse --short HEAD)-SNAPSHOT
LDFLAGS := 
LDFLAGS := $(LDFLAGS) -X main.ProjectName=$(PROJECT_NAME)
LDFLAGS := $(LDFLAGS) -X main.BuildDate=$(BUILD_DATE)
LDFLAGS := $(LDFLAGS) -X main.BuildVersion=$(BUILD_VERSION)

all: test build

build:
	mkdir -p $(DIR_DIST)/bin
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -tags=prod -v $(DIR_SOURCE)/main.go

test:
	mkdir -p $(DIR_DIST)
ifeq ($(OUTPUT),json)
	$(GOTEST) -v ./...  -cover -coverprofile $(DIR_DIST)/cover.out -json > $(DIR_DIST)/test.json
else
	$(GOTEST) -v ./...  -cover
endif

clean:
	#$(GOCLEAN)
	rm -rf $(DIR_OUT)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	echo test
	#$(GOGET) -d -v ./...
