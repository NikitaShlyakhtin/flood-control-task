# Include variables from the .envrc file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/fc: run the cmd/floodcontrol application
.PHONY: run/fc
run/fc:
	go run ./cmd/floodcontrol

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies, format, and vet all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/fc: build the cmd/fc application
.PHONY: build/fc
build/fc:
	@echo 'Building cmd/floodcontrol...'
	go build -o=./bin/floodcontrol ./cmd/floodcontrol
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/floodcontrol ./cmd/floodcontrol