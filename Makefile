# Makefile for Go project

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME_WIN = chabo-api.exe
BINARY_NAME_UNIX = chabo-api.o

# Docker parameters
DOCKERCMD = docker
DOCKERUP = $(DOCKERCMD) compose up
DOCKERBUILD = $(DOCKERCMD) compose build
DOCKERDOWN = $(DOCKERCMD) compose down

# Main build target
all: deps test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME_WIN) -v

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_UNIX) -v

# Clean the build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME_WIN)
	rm -f $(BINARY_NAME_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Install project dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Run the application
run:
	$(DOCKERBUILD) && $(DOCKERUP)

# Format the code
fmt:
	$(GOCMD) fmt ./...

# Lint the code using a linter tool
lint:
	golangci-lint run

# Generate code coverage report
coverage:
	$(GOTEST) -coverprofile='coverage.out' ./...
	$(GOCMD) tool cover -html=coverage.out

# Generate documentation using tools like godoc
doc:
	godoc -http=:6060

# Generate Swagger config
swag:
	swag init -d ./internal/api,./ -g router.go

# Perform a full code quality check (lint, tests, coverage)
check: lint test coverage

.PHONY: all build build-linux clean test deps run fmt lint coverage doc check