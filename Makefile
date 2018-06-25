# Makefile to generate streamline-cli binaries

# Build/version information
NAME    :=$(shell basename `git rev-parse --show-toplevel`)
RELEASE :=$(shell git rev-parse --verify --short HEAD)
VERSION  =$(shell cat VERSION)
BUILD    = $(VERSION)-$(RELEASE)
LDFLAGS  = "-X main.buildVersion=$(BUILD)"

build: 
	GOOS=linux go build -ldflags ${LDFLAGS} -o kissyface-linux
	GOOS=darwin go build -ldflags ${LDFLAGS} -o kissyface-osx
	GOOS=windows go build -ldflags ${LDFLAGS} -o kissyface.exe
