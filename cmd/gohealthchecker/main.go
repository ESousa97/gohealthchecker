package main

import (
	"fmt"
	"os"
	"time"

	"gohealthchecker/internal/checker"
	"gohealthchecker/internal/notifier"
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

	// Initialize the Notifier
	// You can pass the webhook URL via the environment variable WEBHOOK_URL
	webhookURL := os.Getenv("WEBHOOK_URL")
	var n notifier.Notifier

	if webhookURL != "" {
		n = &notifier.WebhookNotifier{URL: webhookURL}
		fmt.Println("Webhook notifier enabled.")
	} else {
		// Use ConsoleNotifier as a fallback for demonstration purposes
		n = &notifier.ConsoleNotifier{}
		fmt.Println("Webhook URL not set. Using ConsoleNotifier. Set WEBHOOK_URL environment variable to enable webhooks.")
	}

	// Initialize the checker service.
	c := checker.Checker{
		Targets:  targets,
		Notifier: n,
	}

	// Start checking the targets every 10 seconds.
	checkInterval := 10 * time.Second
	fmt.Printf("Starting checks every %s...\n", checkInterval)
	c.Start(checkInterval)
}
