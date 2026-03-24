package checker

import (
	"fmt"
	"net/http"
	"time"
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
	Targets []Target
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
	for res := range results {
		timestamp := time.Now().Format(time.RFC3339)
		if res.Error != nil {
			fmt.Printf("[%s] [FAIL] %s - Error: %v - Response Time: %v\n", timestamp, res.Target.URL, res.Error, res.Duration)
		} else if res.Status == http.StatusOK {
			fmt.Printf("[%s] [OK]   %s - Status: %d - Response Time: %v\n", timestamp, res.Target.URL, res.Status, res.Duration)
		} else {
			fmt.Printf("[%s] [WARN] %s - Status: %d - Response Time: %v\n", timestamp, res.Target.URL, res.Status, res.Duration)
		}
	}
}
