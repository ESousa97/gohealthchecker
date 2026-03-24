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
	if webhookURL == "" {
		// Use a public test webhook as fallback for demonstration purposes
		webhookURL = "https://webhook.site/b739e3d1-fb2c-4a04-a6b3-72c5e3c09b96"
		fmt.Println("Webhook URL not set. Using public test webhook:", webhookURL)
	}

	var n notifier.Notifier
	n = &notifier.WebhookNotifier{URL: webhookURL}
	fmt.Println("Webhook notifier enabled.")

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
