GOCC=go
MKDIR=mkdir -p 

BIN_PATH=./build
BIN=./build/eth-light-crawler

.PHONY: build install clean

build:
	$(GOCC) build -o $(BIN)

install:
	$(GOCC) install 

clean:
	rm -r $(BIN_PATH)
