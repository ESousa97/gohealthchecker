package checker

import (
	"fmt"
	"net/http"
	"time"

	"gohealthchecker/internal/notifier"
)

// Target represents an endpoint to be checked.
type Target struct {
	URL string
}

// Result represents the outcome of a single health check.
type Result struct {
	Target   Target
	Status   int
	Duration time.Duration
	Error    error
}

// Checker handles the health checking logic for a list of targets.
type Checker struct {
	Targets  []Target
	Notifier notifier.Notifier
}

// targetState holds the state for a single target to manage alerting and avoid spam.
type targetState struct {
	consecutiveFailures int
	alertSent           bool
}

// Start begins the health check process. It launches one monitoring goroutine
// per target and a central worker to process the results.
func (c *Checker) Start(interval time.Duration) {
	// Central channel for collecting results
	results := make(chan Result)

	// Start the central result worker to format logs
	go c.resultWorker(results)

	// Start a monitoring goroutine for each target
	for _, target := range c.Targets {
		go c.monitorTarget(target, interval, results)
	}

	// Block the main goroutine to keep the application running indefinitely
	select {}
}

// monitorTarget continuously checks a single target at the specified interval.
func (c *Checker) monitorTarget(target Target, interval time.Duration, results chan<- Result) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Perform an initial check immediately
	c.checkTarget(target, results)

	for range ticker.C {
		c.checkTarget(target, results)
	}
}

// checkTarget makes an HTTP GET request to a single target and sends the outcome to the results channel.
func (c *Checker) checkTarget(target Target, results chan<- Result) {
	start := time.Now()

	// HTTP client with a strict timeout to prevent hanging goroutines (network timeouts handling)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(target.URL)
	duration := time.Since(start)

	result := Result{
		Target:   target,
		Duration: duration,
		Error:    err,
	}

	if err == nil {
		result.Status = resp.StatusCode
		// Always close the response body to avoid resource leaks
		resp.Body.Close()
	}

	// Send the result to the central channel
	results <- result
}

// resultWorker reads from the results channel and formats the logs to the terminal.
func (c *Checker) resultWorker(results <-chan Result) {
	// Simple memory state map to track failures per URL
	stateMap := make(map[string]*targetState)

	for res := range results {
		timestamp := time.Now().Format(time.RFC3339)
		url := res.Target.URL

		if _, exists := stateMap[url]; !exists {
			stateMap[url] = &targetState{}
		}
		state := stateMap[url]

		// Consider error or any status other than 200 OK as a failure
		isFailure := res.Error != nil || res.Status != http.StatusOK

		if isFailure {
			if res.Error != nil {
				fmt.Printf("[%s] [FAIL] %s - Error: %v - Response Time: %v\n", timestamp, url, res.Error, res.Duration)
			} else {
				fmt.Printf("[%s] [WARN] %s - Status: %d - Response Time: %v\n", timestamp, url, res.Status, res.Duration)
			}

			state.consecutiveFailures++

			// Trigger alert if it failed twice consecutively and alert hasn't been sent yet
			if state.consecutiveFailures >= 2 && !state.alertSent {
				alertMsg := fmt.Sprintf("🚨 *ALERT*: Service %s is DOWN (Consecutive Failures: %d)", url, state.consecutiveFailures)

				if c.Notifier != nil {
					if err := c.Notifier.Notify(alertMsg); err != nil {
						fmt.Printf("[%s] [NOTIFY ERROR] Failed to send alert for %s: %v\n", timestamp, url, err)
					} else {
						fmt.Printf("[%s] [NOTIFIED] Alert sent for %s\n", timestamp, url)
					}
				}
				state.alertSent = true
			}
		} else {
			fmt.Printf("[%s] [OK]   %s - Status: %d - Response Time: %v\n", timestamp, url, res.Status, res.Duration)

			// Recover service state if it was previously alerted
			if state.alertSent {
				recoveryMsg := fmt.Sprintf("✅ *RECOVERY*: Service %s is UP again.", url)
				if c.Notifier != nil {
					if err := c.Notifier.Notify(recoveryMsg); err != nil {
						fmt.Printf("[%s] [NOTIFY ERROR] Failed to send recovery for %s: %v\n", timestamp, url, err)
					} else {
						fmt.Printf("[%s] [NOTIFIED] Recovery sent for %s\n", timestamp, url)
					}
				}
			}

			// Reset counters on success
			state.consecutiveFailures = 0
			state.alertSent = false
		}
	}
}
