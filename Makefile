BINARY_NAME=sampel_server
OUTPUT_DIR=bin
MAIN_FILE=sample_app/main.go
# Default target: Build the binary
build:
	@echo "Building the binary..."
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Clean target: Remove the binary
clean:
	@echo "Cleaning up..."
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME)

# Install target: Build and move the binary to the output path
install: build
	@echo "Binary installed to $(OUTPUT_DIR)/$(BINARY_NAME)"

# Cross-compile for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_FILE)
