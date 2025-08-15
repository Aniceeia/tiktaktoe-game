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

print_status "Starting cleanup process..."

# Stop and remove containers
print_status "Stopping and removing containers..."
cd docker
docker-compose down --remove-orphans --volumes || true
cd ..

# Remove Docker images
print_status "Removing Docker images..."
docker rmi game:latest || true

# Clean build artifacts
print_status "Cleaning build artifacts..."
rm -rf build/
rm -f game

# Clean Go cache
print_status "Cleaning Go cache..."
go clean -cache -modcache -testcache

# Clean Docker system
print_status "Cleaning Docker system..."
docker system prune -f

print_status "Cleanup completed successfully!"
