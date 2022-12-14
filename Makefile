GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# remove debugging symbols
LDFLAGS += -s -w

NAME=SMBIOSKeygen

PLATFORM := $(shell uname -s)

.PHONY: mac linux native windows

all: native mac linux

native:
	@echo ">  Building native version..." 
ifeq (${PLATFORM},Darwin)
	@$(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME) -tags iokit
else ifeq (${PLATFORM}, Linux)
	@$(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME)
endif
	@echo ">  Done..." 

mac:
	@echo ">  Building Mac versions..." 
	@env GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME)_mac
	@env GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME)_macarm64
	@echo ">  Done..." 

linux:
	@echo ">  Building Linux version..." 
	@env GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME)_linux
	@echo ">  Done..." 

windows:
	@echo ">  Building Windows version..." 
	@env GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="${LDFLAGS}" -o $(NAME)_win
	@echo ">  Done..." 

test:
	@$(GOTEST) .

clean: 
	$(GOCLEAN)
	rm -f $(NAME)
	rm -f $(NAME)_mac
	rm -f $(NAME)_macarm64
	rm -f $(NAME)_linux
