.PHONY: all clean rebuild check config install uninstall test

INSTALL_DIR := $(HOME)/.local/bin
CONFIG_DIR := $(HOME)/.config/cdsl
CONFIG_FILE := $(CONFIG_DIR)/config.toml
EXAMPLE_CONFIG := example-config.toml
BINARY := cdsl

all: $(BINARY) config

check:
	@command -v go >/dev/null 2>&1 || { \
		echo "Error: Go is not installed."; \
		echo "Install it with: sudo apt update && sudo apt install golang"; \
		exit 1; \
	}

	@command -v git >/dev/null 2>&1 || { \
		echo "Error: Git is not installed."; \
		echo "Install it with: sudo apt update && sudo apt install git"; \
		exit 1; \
	}

$(BINARY): check
	@echo "Building CDSL..."
	go build -o $(BINARY) ./cmd/cdsl

config:
	@echo "Setting up CDSL configuration..."
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_FILE) ]; then \
		cp $(EXAMPLE_CONFIG) $(CONFIG_FILE); \
		echo "Created $(CONFIG_FILE)"; \
	else \
		echo "$(CONFIG_FILE) already exists; leaving it unchanged."; \
	fi

clean:
	@echo "Cleaning build files..."
	rm -f $(BINARY)

rebuild: clean
	$(MAKE)

install: $(BINARY)
	@echo "Installing CDSL to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@install -m 755 $(BINARY) $(INSTALL_DIR)/$(BINARY)

uninstall:
	@echo "Removing CDSL from $(INSTALL_DIR)..."
	@rm -f $(INSTALL_DIR)/$(BINARY)
