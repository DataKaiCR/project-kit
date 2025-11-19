.PHONY: build install uninstall clean test help

BINARY_NAME=pk
INSTALL_PATH=/usr/local/bin
MAN_PATH=/usr/local/share/man/man1
BUILD_DIR=bin

help:
	@echo "PK (Project Kit) - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build      Build the binary"
	@echo "  make install    Install binary, man page, and completions"
	@echo "  make uninstall  Remove installed files"
	@echo "  make clean      Remove build artifacts"
	@echo "  make test       Run tests"
	@echo ""
	@echo "After installation:"
	@echo "  pk --help       View help"
	@echo "  man pk          View manual"

build:
	@echo "Building $(BINARY_NAME)..."
	@./build.sh

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✓ Binary installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo ""
	@echo "Installing man page..."
	@sudo mkdir -p $(MAN_PATH)
	@sudo cp docs/pk.1 $(MAN_PATH)/pk.1
	@sudo chmod 644 $(MAN_PATH)/pk.1
	@echo "✓ Man page installed to $(MAN_PATH)/pk.1"
	@echo ""
	@echo "Run './bin/pk install' to complete setup with shell completions"
	@echo "Or reload your shell and run 'pk install'"

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo rm -f $(MAN_PATH)/pk.1
	@echo "✓ Uninstalled"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean"

test:
	@echo "Running tests..."
	@go test ./... -v
