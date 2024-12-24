.PHONY: build clean test

# Variables
BINARY_NAME=monit
SRC_DIR=./main.go
OUTPUT_DIR=./bin


build:
	@echo "Building the project..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"


test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "All tests passed!"


clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleanup complete!"
