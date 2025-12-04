#!/usr/bin/env pwsh
# Test runner script for Windows PowerShell
# Usage: .\scripts\test.ps1 [unit|integration|coverage|all]

param(
    [ValidateSet('unit', 'integration', 'coverage', 'all')]
    [string]$Command = 'all'
)

$ErrorActionPreference = 'Stop'

# Functions
function Write-Header {
    param([string]$Message)
    Write-Host "========================================" -ForegroundColor Blue
    Write-Host $Message -ForegroundColor Blue
    Write-Host "========================================" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Failure {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "! $Message" -ForegroundColor Yellow
}

# Check if PostgreSQL is running
function Test-PostgreSQL {
    Write-Header "Checking PostgreSQL..."
    try {
        $null = Test-NetConnection -ComputerName localhost -Port 5432 -WarningAction SilentlyContinue
        Write-Success "PostgreSQL is running"
        return $true
    }
    catch {
        Write-Failure "PostgreSQL is not running on localhost:5432"
        Write-Warning "Please start PostgreSQL: docker-compose up -d postgres"
        return $false
    }
}

# Run unit tests
function Invoke-UnitTests {
    Write-Header "Running Unit Tests"
    go test -v -race -short .\internal\...\domain\... .\internal\...\application\... .\pkg\...
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Unit tests passed"
    }
    else {
        Write-Failure "Unit tests failed"
        exit 1
    }
}

# Run integration tests
function Invoke-IntegrationTests {
    Write-Header "Running Integration Tests"
    if (-not (Test-PostgreSQL)) {
        exit 1
    }
    
    go test -v -race -run Integration .\tests\integration\...
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Integration tests passed"
    }
    else {
        Write-Failure "Integration tests failed"
        exit 1
    }
}

# Run tests with coverage
function Invoke-Coverage {
    Write-Header "Running Tests with Coverage"
    if (-not (Test-PostgreSQL)) {
        exit 1
    }
    
    go test -v -race -coverprofile=coverage.out -covermode=atomic .\...
    
    Write-Host ""
    Write-Header "Coverage Summary"
    go tool cover -func=coverage.out | Select-Object -Last 1
    
    Write-Host ""
    Write-Header "Generating HTML Coverage Report"
    go tool cover -html=coverage.out -o coverage.html
    Write-Success "Coverage report generated: coverage.html"
    
    # Check coverage threshold
    $coverageLine = go tool cover -func=coverage.out | Select-Object -Last 1
    $coverage = [regex]::Match($coverageLine, '(\d+\.\d+)%').Groups[1].Value
    $threshold = 70.0
    
    if ([double]$coverage -ge $threshold) {
        Write-Success "Coverage $coverage% meets threshold of $threshold%"
    }
    else {
        Write-Failure "Coverage $coverage% is below threshold of $threshold%"
        exit 1
    }
}

# Run all tests
function Invoke-AllTests {
    Write-Header "Running All Tests"
    if (-not (Test-PostgreSQL)) {
        exit 1
    }
    
    go test -v -race .\...
    if ($LASTEXITCODE -eq 0) {
        Write-Success "All tests passed"
    }
    else {
        Write-Failure "Some tests failed"
        exit 1
    }
}

# Main
switch ($Command) {
    'unit' {
        Invoke-UnitTests
    }
    'integration' {
        Invoke-IntegrationTests
    }
    'coverage' {
        Invoke-Coverage
    }
    'all' {
        Invoke-AllTests
    }
}
