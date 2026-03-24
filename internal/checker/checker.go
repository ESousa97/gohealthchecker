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

// Checker handles the health checking logic for a list of targets.
type Checker struct {
	Targets []Target
}

// Start begins checking the targets at the specified interval using a time.Ticker.
func (c *Checker) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Perform an initial check immediately before entering the loop.
	c.checkAll()

	for range ticker.C {
		c.checkAll()
	}
}

// checkAll iterates through all targets and checks them.
func (c *Checker) checkAll() {
	fmt.Printf("\n--- Starting Health Check at %v ---\n", time.Now().Format(time.RFC3339))
	for _, target := range c.Targets {
		c.checkTarget(target)
	}
}

// checkTarget makes an HTTP GET request to a single target and logs the outcome.
func (c *Checker) checkTarget(target Target) {
	start := time.Now()
	
	// Create an HTTP client with a short timeout to prevent hanging.
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(target.URL)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("[FAIL] %s - Error: %v - Response Time: %v\n", target.URL, err, duration)
		return
	}
	// Always close the response body to avoid resource leaks.
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("[OK]   %s - Status: %d - Response Time: %v\n", target.URL, resp.StatusCode, duration)
	} else {
		fmt.Printf("[WARN] %s - Status: %d - Response Time: %v\n", target.URL, resp.StatusCode, duration)
	}
}
