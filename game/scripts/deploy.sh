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

print_status "Starting deployment process..."

# Stop existing containers
print_status "Stopping existing containers..."
cd docker
docker-compose down --remove-orphans || true
cd ..

# Build and start services
print_status "Building and starting services..."
cd docker
docker-compose up --build -d

# Wait for services to be ready
print_status "Waiting for services to be ready..."
sleep 10

# Check if services are healthy
print_status "Checking service health..."
if docker-compose ps | grep -q "healthy"; then
    print_status "Services are healthy!"
else
    print_warning "Some services may not be fully ready yet."
fi

cd ..

print_status "Deployment completed!"
print_status "Application is available at: http://localhost:8080"
print_status "Database is available at: localhost:5432"
print_status "Use 'docker-compose logs -f' to view logs"
