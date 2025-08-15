#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the project root
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run this script from the project root."
    exit 1
fi

print_status "Starting build process..."

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf build/
mkdir -p build/

# Update dependencies
print_status "Updating dependencies..."
go mod tidy
go mod download

# Build for different platforms
print_status "Building for Linux..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/game-linux-amd64 ./cmd/main.go

print_status "Building for macOS..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/game-darwin-amd64 ./cmd/main.go

print_status "Building for Windows..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o build/game-windows-amd64.exe ./cmd/main.go

# Build Docker image
print_status "Building Docker image..."
docker build -f docker/Dockerfile -t game:latest .

print_status "Build completed successfully!"
print_status "Binaries are available in the build/ directory"
print_status "Docker image: game:latest"
