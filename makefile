# Name of the output binary
BINARY_NAME = cheap-fuel 
SRC_DIR = src

.PHONY: all build run test clean

# Default target
all: build

# Build the binary
build:
	go build -o $(BINARY_NAME) ./$(SRC_DIR)

# Run the app
run:
	go run ./$(SRC_DIR)

# Run tests (if you have _test.go files)
test:
	go test ./$(SRC_DIR)/...

# Remove binary
clean:
	rm -f $(BINARY_NAME)
