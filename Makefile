# Makefile to generate streamline-cli binaries

# Build/version information
NAME    :=$(shell basename `git rev-parse --show-toplevel`)
RELEASE :=$(shell git rev-parse --verify --short HEAD)
VERSION  = 0.0.1
BUILD    = $(VERSION)-$(RELEASE)
LDFLAGS  = "-X main.buildVersion=$(BUILD)"

build: 
	go build -ldflags ${LDFLAGS}
