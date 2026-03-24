# GoHealthChecker

A concurrent, high-performance HTTP health checker written in Go. It monitors a list of URLs and sends real-time webhook alerts (e.g., to Slack or Discord) when services go down or recover.

## Features
- **TUI Dashboard:** A real-time terminal user interface (built with Bubble Tea) to monitor service status, latency, and check times visually.
- **Concurrent Checks:** Uses Goroutines to monitor multiple endpoints simultaneously without blocking.
- **Smart Alerting:** Only alerts after 2 consecutive failures to avoid false positives/network blips.
- **Resilience & Retries:** Implements an internal 3-retry mechanism (with a 2-second delay) before a service is considered "failed" to bypass temporary network instabilities.
- **Spam Prevention:** Maintains state in memory. If a service stays down, it won't spam your webhook.
- **Recovery Notifications:** Automatically sends a "Recovery" alert when a failing service returns a `200 OK` status.
- **Dynamic Configuration:** Reads the list of URLs and Webhook credentials dynamically from a `config.yaml` using Viper.

---

## 🚀 How to Run and Test Manually

### Step 1: Run the Application
Open your terminal and run:
```bash
go run cmd/gohealthchecker/main.go
```

### Step 2: Observe the Dashboard
The application will launch a visual dashboard in your terminal:
- **Service URL:** The endpoint being monitored.
- **Status:** Shows `✅ OK` (200 OK) or `❌ FAIL` (error/non-200).
- **Latency:** The time it took to get the response.
- **Last Check:** Timestamp of the last health check.

You can navigate through the rows or press **'q'** to exit.


### Step 3: View the Alerts Live
You can watch the webhook JSON payloads arriving in real-time on our temporary public dashboard:
👉 **[Click here to view live Webhook Alerts on Webhook.site](https://webhook.site/#!/view/b739e3d1-fb2c-4a04-a6b3-72c5e3c09b96)**

*(Note: If you want to use your own Discord or Slack webhook, simply pass it via the `WEBHOOK_URL` environment variable before running the command).*

---

## 🧪 Running Automated Tests

The project includes robust automated tests that utilize Go's `httptest` package to simulate servers and webhook endpoints natively, ensuring everything works without relying on external internet connections.

Run all tests with:
```bash
go test -v ./...
```

You should see an output indicating that the webhook notification logic and the state management (spam prevention) passed successfully.
