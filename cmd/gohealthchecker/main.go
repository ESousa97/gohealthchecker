package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"gohealthchecker/internal/checker"
	"gohealthchecker/internal/notifier"
)

func main() {
	fmt.Println("Initializing GoHealthChecker...")

	// Configure Viper to read config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Default values
	viper.SetDefault("check_interval", "10s")
	viper.SetDefault("targets", []string{})
	viper.SetDefault("webhook_url", "")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	targetURLs := viper.GetStringSlice("targets")
	var targets []checker.Target
	for _, url := range targetURLs {
		targets = append(targets, checker.Target{URL: url})
	}

	if len(targets) == 0 {
		log.Fatal("No targets configured to monitor. Please add targets to config.yaml.")
	}

	webhookURL := viper.GetString("webhook_url")
	var n notifier.Notifier

	if webhookURL != "" {
		n = &notifier.WebhookNotifier{URL: webhookURL}
		fmt.Println("Webhook notifier enabled via config.")
	} else {
		// Use ConsoleNotifier as a fallback
		n = &notifier.ConsoleNotifier{}
		fmt.Println("Webhook URL not set in config. Using ConsoleNotifier.")
	}

	// Initialize the checker service
	c := checker.Checker{
		Targets:  targets,
		Notifier: n,
	}

	// Start checking the targets based on the configured interval
	checkInterval := viper.GetDuration("check_interval")
	fmt.Printf("Starting checks every %s...\n", checkInterval)
	c.Start(checkInterval)
}
