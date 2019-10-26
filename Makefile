.PHONY: all check-path build clean install uninstall fmt simplify run
.DEFAULT_GOAL: $(BIN_FILE)

SHELL := /bin/bash

# export GOPATH := $(shell pwd)

export GO111MODULE=on

PROJECT_NAME = authz

BUILD_DIR = build
BIN_DIR = bin
BIN_FILE = $(PROJECT_NAME)

SRC_DIR = ./
SRC_PKGS = $(shell GOPATH=$(GOPATH); go list $(SRC_DIR)/...)
SRC_FILES = $(shell find . -type f -name '*.go' -path "$(SRC_DIR)/*")

# Get version constant
VERSION := $(shell cat $(SRC_DIR)/main.go | grep "const version = " | awk '{print $$NF}' | sed -e 's/^.//' -e 's/.$$//')
BUILD := $(shell git rev-parse HEAD)

# Use linker flags to provide version/build settings to the binary
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

DOCKER_COMPOSE_CMD = docker-compose -p $(PROJECT_NAME) -f docker/docker-compose.yml


all: dep build

dep:
	@echo "Downloading dependencies..."
	@export GO111MODULE=on
	@go mod tidy
	@go mod download
	@go mod vendor
	@echo "Finish..."

build: dep
	@echo "Compiling $(BUILD_DIR)/$(BIN_FILE)..."
	@go build $(LDFLAGS) -i -o $(BUILD_DIR)/$(BIN_FILE) $(SRC_DIR)/main.go
	@echo "Finish..."

clean:
	@rm -rf $(BIN_DIR)/$(BIN_FILE)
	@rm -rf $(BUILD_DIR)
	@rm -rf pkg/
	@rm -rf bin/
	@rm -rf templates/
	@rm -rf default/

install:
	@go install $(LDFLAGS) $(SRC_PKGS)
	@cp $(BIN_DIR)/$(BIN_FILE) /usr/local/bin/
# 	@cp etc/config.yml /etc/microservice-email.yml
	@echo "Instalation complete..."

uninstall:
	@rm -f $$(which $(BIN_FILE))
# 	@rm -f /etc/microservice-email.yml

run: build
	$(BUILD_DIR)/$(BIN_FILE) consume

# -------------------------------------------------------------------
# -								Docker								-
# -------------------------------------------------------------------

docker_build:
	@$(DOCKER_COMPOSE_CMD) build $(PROJECT_NAME)

docker_shell:
	@$(DOCKER_COMPOSE_CMD) run --rm $(PROJECT_NAME) /bin/bash

docker_run:
	@$(DOCKER_COMPOSE_CMD) run --rm --service-ports --name $(PROJECT_NAME) $(PROJECT_NAME) \
		/bin/bash -ci "make dep && make run"

up: ## Start the container
	@echo "Starting Container"
	@$(DOCKER_COMMAND) up -d --build

down: ## Bring Down the container
	@echo "Stopping Container"
	@$(DOCKER_COMMAND) down

clean: ## Remove the container with volume
	@echo "Removing Container"
	@$(DOCKER_COMMAND) down -v