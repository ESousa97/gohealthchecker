// Package checker provides the core logic for concurrent HTTP health checks,
// managing target states, retries, and alerting triggers.
package checker

import (
	"fmt"
	"net/http"
	"time"

	"gohealthchecker/internal/notifier"
)

// Target represents an endpoint to be monitored for health.
type Target struct {
	// URL is the full address of the endpoint (e.g., "https://api.example.com").
	URL string
}

// Result represents the outcome of a single health check execution.
type Result struct {
	Target    Target        // The target that was checked.
	Status    int           // HTTP status code returned (e.g., 200).
	Duration  time.Duration // Time taken to receive the response.
	Error     error         // Any network or protocol error encountered.
	LastCheck time.Time     // Timestamp of when the check was performed.
}

// Checker manages the orchestration of health checks for a set of [Target]s.
// It uses a [notifier.Notifier] to send alerts based on check outcomes.
type Checker struct {
	Targets  []Target
	Notifier notifier.Notifier
}

// TargetState holds the state for a single target to manage alerting and avoid spam.
type TargetState struct {
	ConsecutiveFailures int
	AlertSent           bool
}

// Start begins the health check process in the background.
// It launches one monitoring goroutine per target and returns a channel
// through which [Result]s are asynchronously delivered.
// The checks repeat at the specified interval.
func (c *Checker) Start(interval time.Duration) <-chan Result {
	// Central channel for collecting results
	results := make(chan Result)

	// Start a monitoring goroutine for each target
	for _, target := range c.Targets {
		go c.monitorTarget(target, interval, results)
	}

	return results
}

// RunWorker starts the central result worker to format logs to console and handle alerts.
// This is used for non-TUI mode.
func (c *Checker) RunWorker(results <-chan Result) {
	stateMap := make(map[string]*TargetState)
	for res := range results {
		timestamp := res.LastCheck.Format(time.RFC3339)
		url := res.Target.URL

		// Process result for alerts
		c.ProcessResult(res, stateMap)

		if res.Error != nil {
			fmt.Printf("[%s] [FAIL] %s - Error: %v - Response Time: %v\n", timestamp, url, res.Error, res.Duration)
		} else if res.Status == http.StatusOK {
			fmt.Printf("[%s] [OK]   %s - Status: %d - Response Time: %v\n", timestamp, url, res.Status, res.Duration)
		} else {
			fmt.Printf("[%s] [WARN] %s - Status: %d - Response Time: %v\n", timestamp, url, res.Status, res.Duration)
		}
	}
}

// ProcessResult evaluates a [Result] against the current [TargetState] to determine
// if an alert or a recovery notification should be triggered.
//
// It returns isAlerted=true if a new failure alert was sent,
// and isRecovered=true if a recovery notification was sent.
func (c *Checker) ProcessResult(res Result, stateMap map[string]*TargetState) (isAlerted bool, isRecovered bool) {
	url := res.Target.URL
	if _, exists := stateMap[url]; !exists {
		stateMap[url] = &TargetState{}
	}
	state := stateMap[url]

	isFailure := res.Error != nil || res.Status != http.StatusOK

	if isFailure {
		state.ConsecutiveFailures++
		if state.ConsecutiveFailures >= 2 && !state.AlertSent {
			state.AlertSent = true
			alertMsg := fmt.Sprintf("🚨 *ALERT*: Service %s is DOWN (Consecutive Failures: %d)", url, state.ConsecutiveFailures)
			if c.Notifier != nil {
				c.Notifier.Notify(alertMsg)
			}
			return true, false
		}
	} else {
		if state.AlertSent {
			state.AlertSent = false
			state.ConsecutiveFailures = 0
			recoveryMsg := fmt.Sprintf("✅ *RECOVERY*: Service %s is UP again.", url)
			if c.Notifier != nil {
				c.Notifier.Notify(recoveryMsg)
			}
			return false, true
		}
		state.ConsecutiveFailures = 0
		state.AlertSent = false
	}
	return false, false
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
// Implements retry logic: if a check fails, it retries up to 3 times with a 2-second delay.
func (c *Checker) checkTarget(target Target, results chan<- Result) {
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	var lastErr error
	var lastStatus int
	var duration time.Duration
	lastCheck := time.Now()

	for attempt := 0; attempt <= maxRetries; attempt++ {
		start := time.Now()

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Get(target.URL)
		duration = time.Since(start)

		lastErr = err
		if err == nil {
			lastStatus = resp.StatusCode
			resp.Body.Close()

			if lastStatus == http.StatusOK {
				break
			}
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	result := Result{
		Target:    target,
		Status:    lastStatus,
		Duration:  duration,
		Error:     lastErr,
		LastCheck: lastCheck,
	}

	results <- result
}
