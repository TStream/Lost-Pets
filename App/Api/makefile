# GO Build settings 
GOOS?=windows
GOARCH?=amd64
NOW?=$(shell date +%Y-%m-%dT%T%z)
BIN?=$(CURDIR)/bin

# --match limits the tag returned to version tags
# file .version is only used if git describe returns nil
VERSION?=$(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

# Output file location/name
LOSTPETS=$(BIN)/lost-pets
MIGRATIONS=$(BIN)/db-migrations
# Default CMD, run when make is call without a target
.PHONY: default
default: help;

# Build Commands
all: fmt lint | $(MIGRATIONS) $(LOSTPETS) ## Build all cmd binaries

$(MIGRATIONS): ## Build migrations binary
	go build -o $(MIGRATIONS) lostpets/cmd/db

db-migrations: $(MIGRATIONS) ## Build migrations binary

$(LOSTPETS):  ## Build lost-pets server
	go build \
	-ldflags '-X main.version=$(VERSION) -X main.timestamp=$(NOW)' \
	-o $(BIN)/lost-pets lostpets/cmd/server

lost-pets: $(LOSTPETS) ## Build lost-pets server