package main

// Example of how to use build-time version information
// This shows how to embed version info in your binaries

import (
	"fmt"
)

// These variables are set at build time using -ldflags
var (
	Version   = "dev"      // Set via -X main.Version=1.0.0
	BuildTime = "unknown"  // Set via -X main.BuildTime=2024-01-01
	GitCommit = "none"     // Set via -X main.GitCommit=abc1234
)

// PrintVersion prints version information
func PrintVersion() {
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
}

// Example usage in main.go:
//
// func main() {
//     // Add --version flag
//     versionFlag := flag.Bool("version", false, "Print version information")
//     flag.Parse()
//
//     if *versionFlag {
//         PrintVersion()
//         os.Exit(0)
//     }
//
//     // Rest of your application...
// }

// Example build command:
// go build -ldflags="-X main.Version=1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.GitCommit=$(git rev-parse --short HEAD)" -o myapp
