package checker

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

// MockNotifier captures notifications for testing purposes.
type MockNotifier struct {
	Notifications []string
}

func (m *MockNotifier) Notify(message string) error {
	m.Notifications = append(m.Notifications, message)
	return nil
}

func TestResultWorker_StateAndAlerts(t *testing.T) {
	mockNotifier := &MockNotifier{}

	c := &Checker{
		Notifier: mockNotifier,
	}

	resultsChan := make(chan Result)

	// Start the worker in a goroutine
	go c.resultWorker(resultsChan)

	targetURL := "http://test.local"
	target := Target{URL: targetURL}

	// 1st Failure: Should NOT trigger an alert (needs 2 consecutive)
	resultsChan <- Result{Target: target, Status: http.StatusInternalServerError, Error: errors.New("timeout")}
	time.Sleep(50 * time.Millisecond) // Give worker time to process

	if len(mockNotifier.Notifications) != 0 {
		t.Fatalf("Expected 0 notifications after 1st failure, got %d", len(mockNotifier.Notifications))
	}

	// 2nd Failure: Should trigger a DOWN alert
	resultsChan <- Result{Target: target, Status: http.StatusInternalServerError, Error: errors.New("timeout")}
	time.Sleep(50 * time.Millisecond)

	if len(mockNotifier.Notifications) != 1 {
		t.Fatalf("Expected 1 notification after 2nd failure, got %d", len(mockNotifier.Notifications))
	}

	// 3rd Failure: Should NOT trigger a new alert (prevents spam)
	resultsChan <- Result{Target: target, Status: http.StatusInternalServerError, Error: errors.New("timeout")}
	time.Sleep(50 * time.Millisecond)

	if len(mockNotifier.Notifications) != 1 {
		t.Fatalf("Expected still 1 notification after 3rd failure (spam prevention), got %d", len(mockNotifier.Notifications))
	}

	// 1st Success (Recovery): Should trigger an UP/RECOVERY alert
	resultsChan <- Result{Target: target, Status: http.StatusOK, Error: nil}
	time.Sleep(50 * time.Millisecond)

	if len(mockNotifier.Notifications) != 2 {
		t.Fatalf("Expected 2 notifications after recovery, got %d", len(mockNotifier.Notifications))
	}

	// Close the channel to gracefully stop the worker
	close(resultsChan)
}
