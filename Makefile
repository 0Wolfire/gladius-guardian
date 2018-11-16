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

GUARD_SRC=$(SRC_DIR)/main.go
GUARD_DEST=$(DST_DIR)/gladius-guardian$(BINARY_SUFFIX)

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

executable:
	$(GOBUILD) -o $(GUARD_DEST) $(GUARD_SRC)
