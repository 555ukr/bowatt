APP_NAME := bowatt
BUILD_DIR := bin
MAIN := cmd/main.go

.PHONY: build run clean test fmt vet

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

