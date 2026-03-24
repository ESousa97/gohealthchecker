package main

import (
	"fmt"
	"time"

	"gohealthchecker/internal/checker"
)

func main() {
	fmt.Println("Initializing GoHealthChecker...")

	// Define the list of URLs to check via the Target struct.
	targets := []checker.Target{
		{URL: "https://www.google.com"},
		{URL: "https://pkg.go.dev"},
		{URL: "https://httpstat.us/404"},
		{URL: "http://invalid.local.domain"},
	}

	// Initialize the checker service.
	c := checker.Checker{
		Targets: targets,
	}

	// Start checking the targets every 10 seconds.
	checkInterval := 10 * time.Second
	fmt.Printf("Starting checks every %s...\n", checkInterval)
	c.Start(checkInterval)
}
