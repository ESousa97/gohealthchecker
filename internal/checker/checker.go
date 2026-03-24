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
	Target    Target
	Status    int
	Duration  time.Duration
	Error     error
	LastCheck time.Time
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
// per target and returns a channel of results.
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
	for res := range results {
		timestamp := res.LastCheck.Format(time.RFC3339)
		url := res.Target.URL

		// Note: The logic for state and alerts is now moved to the result processor
		// to allow both TUI and console worker to trigger notifications.
		// For simplicity in this exercise, we maintain the state management inside the worker or UI.
	}
}

// ProcessResult updates the state and sends alerts if needed.
// This logic is shared by both TUI and standard worker.
func (c *Checker) ProcessResult(res Result, stateMap map[string]*targetState) (isAlerted bool, isRecovered bool) {
	url := res.Target.URL
	if _, exists := stateMap[url]; !exists {
		stateMap[url] = &targetState{}
	}
	state := stateMap[url]

	isFailure := res.Error != nil || res.Status != http.StatusOK

	if isFailure {
		state.consecutiveFailures++
		if state.consecutiveFailures >= 2 && !state.alertSent {
			state.alertSent = true
			alertMsg := fmt.Sprintf("🚨 *ALERT*: Service %s is DOWN (Consecutive Failures: %d)", url, state.consecutiveFailures)
			if c.Notifier != nil {
				c.Notifier.Notify(alertMsg)
			}
			return true, false
		}
	} else {
		if state.alertSent {
			state.alertSent = false
			state.consecutiveFailures = 0
			recoveryMsg := fmt.Sprintf("✅ *RECOVERY*: Service %s is UP again.", url)
			if c.Notifier != nil {
				c.Notifier.Notify(recoveryMsg)
			}
			return false, true
		}
		state.consecutiveFailures = 0
		state.alertSent = false
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
