#!/bin/bash

# E-Commerce Platform Build Script
# Compiles binaries for multiple platforms: Linux, macOS, Windows

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="ecommerce"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
ADMIN_CMD="cmd/admin/main.go"
STOREFRONT_CMD="cmd/storefront/main.go"

# Build timestamp
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Print colored message
print_msg() {
    echo -e "${2}${1}${NC}"
}

# Print header
print_header() {
    echo ""
    print_msg "════════════════════════════════════════════════════════" "$BLUE"
    print_msg "  $1" "$BLUE"
    print_msg "════════════════════════════════════════════════════════" "$BLUE"
    echo ""
}

# Build for a specific platform
build_binary() {
    local os=$1
    local arch=$2
    local app=$3
    local source=$4
    local output_name=$5

    # Set file extension for Windows
    local ext=""
    if [ "$os" = "windows" ]; then
        ext=".exe"
    fi

    local output_dir="${BUILD_DIR}/${os}-${arch}"
    local output_file="${output_dir}/${output_name}${ext}"

    # Create output directory
    mkdir -p "$output_dir"

    print_msg "  Building ${app} for ${os}/${arch}..." "$YELLOW"

    # Build with ldflags for version info
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
        -o "$output_file" \
        "$source"

    if [ $? -eq 0 ]; then
        local size=$(du -h "$output_file" | cut -f1)
        print_msg "    ✓ Built: ${output_file} (${size})" "$GREEN"
    else
        print_msg "    ✗ Failed to build ${app} for ${os}/${arch}" "$RED"
        return 1
    fi
}

# Build all platforms
build_all() {
    local app=$1
    local source=$2
    local binary_name=$3

    print_header "Building ${app}"

    # Linux builds
    build_binary "linux" "amd64" "$app" "$source" "$binary_name"
    build_binary "linux" "arm64" "$app" "$source" "$binary_name"

    # macOS builds
    build_binary "darwin" "amd64" "$app" "$source" "$binary_name"
    build_binary "darwin" "arm64" "$app" "$source" "$binary_name"

    # Windows builds
    build_binary "windows" "amd64" "$app" "$source" "$binary_name"
}

# Build specific platform
build_platform() {
    local platform=$1
    local app=$2
    local source=$3
    local binary_name=$4

    case $platform in
        linux)
            build_binary "linux" "amd64" "$app" "$source" "$binary_name"
            build_binary "linux" "arm64" "$app" "$source" "$binary_name"
            ;;
        macos|darwin)
            build_binary "darwin" "amd64" "$app" "$source" "$binary_name"
            build_binary "darwin" "arm64" "$app" "$source" "$binary_name"
            ;;
        windows)
            build_binary "windows" "amd64" "$app" "$source" "$binary_name"
            ;;
        *)
            print_msg "Unknown platform: $platform" "$RED"
            print_msg "Valid platforms: linux, macos, windows, all" "$YELLOW"
            exit 1
            ;;
    esac
}

# Clean build directory
clean() {
    print_header "Cleaning build directory"
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        print_msg "  ✓ Removed ${BUILD_DIR}/" "$GREEN"
    else
        print_msg "  Build directory doesn't exist" "$YELLOW"
    fi
}

# Create release archives
create_archives() {
    print_header "Creating release archives"

    cd "$BUILD_DIR"

    for dir in */; do
        dir_name="${dir%/}"
        archive_name="${APP_NAME}-${dir_name}-${VERSION}.tar.gz"

        print_msg "  Creating ${archive_name}..." "$YELLOW"
        tar -czf "$archive_name" "$dir_name"

        if [ $? -eq 0 ]; then
            local size=$(du -h "$archive_name" | cut -f1)
            print_msg "    ✓ Created: ${archive_name} (${size})" "$GREEN"
        fi
    done

    cd ..
}

# Print build summary
print_summary() {
    print_header "Build Summary"

    print_msg "Version:     ${VERSION}" "$BLUE"
    print_msg "Build Time:  ${BUILD_TIME}" "$BLUE"
    print_msg "Git Commit:  ${GIT_COMMIT}" "$BLUE"
    echo ""

    if [ -d "$BUILD_DIR" ]; then
        print_msg "Build artifacts:" "$GREEN"
        find "$BUILD_DIR" -type f \( -name "admin*" -o -name "storefront*" -o -name "*.tar.gz" \) -exec ls -lh {} \; | awk '{print "  " $9 " (" $5 ")"}'
    fi
    echo ""
}

# Show usage
usage() {
    cat << EOF
E-Commerce Platform Build Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    all                 Build for all platforms (Linux, macOS, Windows)
    linux               Build for Linux (amd64, arm64)
    macos|darwin        Build for macOS (amd64, arm64)
    windows             Build for Windows (amd64)
    clean               Clean build directory
    release             Build all platforms and create release archives
    help                Show this help message

Options:
    --admin-only        Build only admin binary
    --storefront-only   Build only storefront binary
    --version VERSION   Set build version (default: 1.0.0)

Environment Variables:
    VERSION            Build version (default: 1.0.0)

Examples:
    # Build for all platforms
    $0 all

    # Build for Linux only
    $0 linux

    # Build for macOS with specific version
    VERSION=2.0.0 $0 macos

    # Build only admin for all platforms
    $0 all --admin-only

    # Create release archives
    $0 release

    # Clean build directory
    $0 clean

EOF
}

# Main script logic
main() {
    local command=${1:-help}
    local admin_only=false
    local storefront_only=false

    # Parse options
    shift || true
    while [[ $# -gt 0 ]]; do
        case $1 in
            --admin-only)
                admin_only=true
                shift
                ;;
            --storefront-only)
                storefront_only=true
                shift
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            *)
                print_msg "Unknown option: $1" "$RED"
                usage
                exit 1
                ;;
        esac
    done

    # Print build info
    print_header "E-Commerce Platform Build"
    print_msg "Version: ${VERSION}" "$BLUE"
    print_msg "Build Time: ${BUILD_TIME}" "$BLUE"
    print_msg "Git Commit: ${GIT_COMMIT}" "$BLUE"

    case $command in
        all)
            if [ "$storefront_only" = false ]; then
                build_all "Admin" "$ADMIN_CMD" "admin"
            fi
            if [ "$admin_only" = false ]; then
                build_all "Storefront" "$STOREFRONT_CMD" "storefront"
            fi
            print_summary
            ;;
        linux|macos|darwin|windows)
            if [ "$storefront_only" = false ]; then
                print_header "Building Admin for $command"
                build_platform "$command" "Admin" "$ADMIN_CMD" "admin"
            fi
            if [ "$admin_only" = false ]; then
                print_header "Building Storefront for $command"
                build_platform "$command" "Storefront" "$STOREFRONT_CMD" "storefront"
            fi
            print_summary
            ;;
        clean)
            clean
            ;;
        release)
            if [ "$storefront_only" = false ]; then
                build_all "Admin" "$ADMIN_CMD" "admin"
            fi
            if [ "$admin_only" = false ]; then
                build_all "Storefront" "$STOREFRONT_CMD" "storefront"
            fi
            create_archives
            print_summary
            ;;
        help|--help|-h)
            usage
            exit 0
            ;;
        *)
            print_msg "Unknown command: $command" "$RED"
            usage
            exit 1
            ;;
    esac

    print_msg "\n✓ Build completed successfully!\n" "$GREEN"
}

# Run main function
main "$@"
