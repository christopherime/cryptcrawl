.PHONY: build run test clean ssh-key ssh-connect lint new-dungeon

# Default target
all: build

# Build the application
build:
	mkdir -p bin
	go build -o bin/cryptcrawl ./cmd/cryptcrawl

# Run the application
run: build
	./bin/cryptcrawl

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f bin/cryptcrawl

# Generate SSH key
ssh-key:
	mkdir -p .ssh
	ssh-keygen -t ed25519 -f .ssh/cryptcrawl_ed25519 -N ""

# Connect to the SSH server
ssh-connect:
	ssh localhost -p 23234

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go get -u github.com/charmbracelet/wish
	go get -u github.com/charmbracelet/bubbletea
	go get -u github.com/charmbracelet/lipgloss
	go get -u github.com/charmbracelet/bubbles

# Create a new dungeon
new-dungeon:
	mkdir -p dungeons
	cp internal/dungeon/examples/forgotten_crypt.json dungeons/my_dungeon.json
	@echo "Created a new dungeon template at dungeons/my_dungeon.json"
	@echo "Edit this file to create your own dungeon!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  ssh-key      - Generate SSH key"
	@echo "  ssh-connect  - Connect to the SSH server"
	@echo "  lint         - Run linter"
	@echo "  deps         - Install dependencies"
	@echo "  new-dungeon  - Create a new dungeon template"
	@echo "  help         - Show this help message"
