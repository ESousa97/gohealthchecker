// Package notifier provides abstractions and implementations for sending
// health check alerts to various destinations like Webhooks or Console.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Notifier defines the common interface for sending alert notifications.
type Notifier interface {
	// Notify sends the given message to the notification destination.
	// It returns an error if the delivery fails.
	Notify(message string) error
}

// WebhookNotifier implements the [Notifier] interface to send JSON payloads
// to a specific URL, compatible with Slack, Discord, and others.
type WebhookNotifier struct {
	// URL is the destination endpoint for the HTTP POST request.
	URL string
}

// WebhookPayload represents the standard JSON structure expected by most
// webhook consumers (like Slack or Discord).
type WebhookPayload struct {
	Content string `json:"content"`
}

// Notify sends a formatted JSON message to the configured webhook URL.
func (n *WebhookNotifier) Notify(message string) error {
	if n.URL == "" {
		return nil // Webhook URL not configured, silently ignore
	}

	payload := WebhookPayload{Content: message}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Short timeout to avoid blocking the caller
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(n.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code from webhook: %d", resp.StatusCode)
	}

	return nil
}

// ConsoleNotifier is a fallback Notifier that prints to the console if no webhook is set.
type ConsoleNotifier struct{}

// Notify prints the message to standard output.
func (c *ConsoleNotifier) Notify(message string) error {
	fmt.Printf("[ALERT] %s\n", message)
	return nil
}
