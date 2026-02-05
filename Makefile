BINARY = cpumon
BUILD_DIR = build

.PHONY: build build-optimized run install clean lint

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) .

build-optimized:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -gcflags="-l=4" -o $(BUILD_DIR)/$(BINARY) .

run: build
	@$(BUILD_DIR)/$(BINARY)

install: build-optimized
	sudo cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

lint:
	go vet ./...
	golangci-lint run
