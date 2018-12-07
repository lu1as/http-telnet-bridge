BIN_DIR=bin
BIN=$(BIN_DIR)/bridge

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN)
