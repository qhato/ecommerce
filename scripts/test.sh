#!/bin/bash

# Test runner script with colored output
# Usage: ./scripts/test.sh [unit|integration|coverage|all]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}! $1${NC}"
}

# Check if PostgreSQL is running
check_postgres() {
    print_header "Checking PostgreSQL..."
    if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
        print_error "PostgreSQL is not running on localhost:5432"
        print_warning "Please start PostgreSQL: docker-compose up -d postgres"
        exit 1
    fi
    print_success "PostgreSQL is running"
}

# Check if Redis is running (optional)
check_redis() {
    if ! redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
        print_warning "Redis is not running (optional for tests)"
    else
        print_success "Redis is running"
    fi
}

# Run unit tests
run_unit_tests() {
    print_header "Running Unit Tests"
    go test -v -race -short ./internal/*/domain/... ./internal/*/application/... ./pkg/...
    print_success "Unit tests passed"
}

# Run integration tests
run_integration_tests() {
    print_header "Running Integration Tests"
    check_postgres
    go test -v -race -run Integration ./tests/integration/...
    print_success "Integration tests passed"
}

# Run all tests with coverage
run_coverage() {
    print_header "Running Tests with Coverage"
    check_postgres
    
    go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
    
    echo ""
    print_header "Coverage Summary"
    go tool cover -func=coverage.out | tail -n 1
    
    echo ""
    print_header "Generating HTML Coverage Report"
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
    
    # Check coverage threshold
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
    threshold=70.0
    
    if (( $(echo "$coverage >= $threshold" | bc -l) )); then
        print_success "Coverage $coverage% meets threshold of $threshold%"
    else
        print_error "Coverage $coverage% is below threshold of $threshold%"
        exit 1
    fi
}

# Run all tests
run_all_tests() {
    print_header "Running All Tests"
    check_postgres
    check_redis
    
    go test -v -race ./...
    print_success "All tests passed"
}

# Main
case "${1:-all}" in
    unit)
        run_unit_tests
        ;;
    integration)
        run_integration_tests
        ;;
    coverage)
        run_coverage
        ;;
    all)
        run_all_tests
        ;;
    *)
        echo "Usage: $0 {unit|integration|coverage|all}"
        echo ""
        echo "Commands:"
        echo "  unit        - Run unit tests only (fast)"
        echo "  integration - Run integration tests (requires PostgreSQL)"
        echo "  coverage    - Run all tests with coverage report"
        echo "  all         - Run all tests (default)"
        exit 1
        ;;
esac
