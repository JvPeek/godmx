BUILD_DIR := build

.PHONY: all linux macos windows clean

all: linux macos windows

linux: $(BUILD_DIR)/godmx_linux

macos: $(BUILD_DIR)/godmx_macos

windows: $(BUILD_DIR)/godmx_windows.exe

$(BUILD_DIR)/godmx_linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ .

$(BUILD_DIR)/godmx_macos:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $@ .

$(BUILD_DIR)/godmx_windows.exe:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $@ .

clean:
	rm -rf $(BUILD_DIR)
