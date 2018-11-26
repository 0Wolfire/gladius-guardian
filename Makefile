##
## Makefile to test and build the gladius binaries
##

##
# GLOBAL VARIABLES
##

# if we are running on a windows machine
# we need to append a .exe to the
# compiled binary
BINARY_SUFFIX=
ifeq ($(OS),Windows_NT)
	BINARY_SUFFIX=.exe
endif

ifeq ($(GOOS),windows)
	BINARY_SUFFIX=.exe
endif

# code source and build directories
SRC_DIR=.
DST_DIR=./build

BINARY=gladius-guardian$(BINARY_SUFFIX)
GUARD_SRC=$(SRC_DIR)/main.go
GUARD_DEST=$(DST_DIR)/$(BINARY)

# commands for go
GOMOD=GO111MODULE=on
GOBUILD=$(GOMOD) go build
GOTEST=$(GOMOD) go test
GOCLEAN=$(GOMOD) go clean

##
# MAKE TARGETS
##

# general make targets
all: 
	make clean
	# make lint
	make executable

clean:
	rm -rf ./build/*
	$(GOMOD) go mod tidy
	$(GOCLEAN)

lint:
	gometalinter --linter='vet:go tool vet -printfuncs=Infof,Debugf,Warningf,Errorf:PATH:LINE:MESSAGE' main.go

test: $(CTL_SRC)
	$(GOTEST) $(GUARD_SRC)

# Made for macOS at the moment
# Install gcc cross compilers for macOS
# `brew install mingw-w64` - windows
# `brew install FiloSottile/musl-cross/musl-cross` - linux
release: clean release-win release-linux release-mac

release-win:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/windows/$(BINARY).exe $(GUARD_SRC)
release-linux:
	CGO_ENABLED=1 CC=x86_64-linux-musl-gcc GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/linux/$(BINARY) $(GUARD_SRC)
release-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DST_DIR)/release/macos/$(BINARY) $(GUARD_SRC)

executable:
	$(GOBUILD) -o $(GUARD_DEST) $(GUARD_SRC)
