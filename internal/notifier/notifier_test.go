package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookNotifier_Notify_Success(t *testing.T) {
	// Create a mock HTTP server that simulates a successful webhook response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		var payload WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if payload.Content != "test message" {
			t.Errorf("Expected content 'test message', got '%s'", payload.Content)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	notifier := &WebhookNotifier{URL: mockServer.URL}

	err := notifier.Notify("test message")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestWebhookNotifier_Notify_Failure(t *testing.T) {
	// Create a mock HTTP server that simulates a failed webhook response (e.g., 500 Internal Server Error)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	notifier := &WebhookNotifier{URL: mockServer.URL}

	err := notifier.Notify("test message")
	if err == nil {
		t.Fatal("Expected an error for status 500, but got nil")
	}
}

func TestWebhookNotifier_Notify_EmptyURL(t *testing.T) {
	notifier := &WebhookNotifier{URL: ""}
	err := notifier.Notify("test message")
	if err != nil {
		t.Fatalf("Expected no error when URL is empty (silent ignore), got: %v", err)
	}
}
