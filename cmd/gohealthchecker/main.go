package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"gohealthchecker/internal/checker"
	"gohealthchecker/internal/notifier"
	"gohealthchecker/internal/tui"
)

func main() {
	// Configure Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	targetURLs := viper.GetStringSlice("targets")
	var targets []checker.Target
	for _, url := range targetURLs {
		targets = append(targets, checker.Target{URL: url})
	}

	webhookURL := viper.GetString("webhook_url")
	var n notifier.Notifier
	if webhookURL != "" {
		n = &notifier.WebhookNotifier{URL: webhookURL}
	} else {
		n = &notifier.ConsoleNotifier{}
	}

	// Initialize checker
	c := checker.Checker{
		Targets:  targets,
		Notifier: n,
	}

	// Start checking (returns a channel of results)
	checkInterval := viper.GetDuration("check_interval")
	results := c.Start(checkInterval)

	// Check if TUI should be enabled (default yes if terminal is interactive)
	// For simplicity, we'll always run the TUI in this phase
	fmt.Println("Launching TUI Monitoring Dashboard...")
	if err := tui.StartUI(&c, results); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
