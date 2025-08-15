#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[PASS]${NC} $1"
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

print_status "Starting test process..."

# Check if application is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    print_error "Application is not running. Please start it first."
    print_status "You can start it with: ./scripts/deploy.sh"
    exit 1
fi

print_status "Application is running. Starting tests..."

# Run integration tests
print_status "Running integration tests..."
chmod +x tests/test_api.sh
chmod +x tests/test_all.sh

print_status "Running API tests..."
./tests/test_api.sh
print_status "API tests completed successfully!"

print_status "Running full integration tests..."
./tests/test_all.sh
print_status "full integration tests completed successfully!"


print_status "All tests completed successfully!"
