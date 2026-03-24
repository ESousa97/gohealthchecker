# GoHealthChecker

A concurrent, high-performance HTTP health checker written in Go. It monitors a list of URLs and sends real-time webhook alerts (e.g., to Slack or Discord) when services go down or recover.

## Features
- **Concurrent Checks:** Uses Goroutines to monitor multiple endpoints simultaneously without blocking.
- **Smart Alerting:** Only alerts after 2 consecutive failures to avoid false positives/network blips.
- **Spam Prevention:** Maintains state in memory. If a service stays down, it won't spam your webhook.
- **Recovery Notifications:** Automatically sends a "Recovery" alert when a failing service returns a `200 OK` status.
- **Fallback Console Logs:** If no webhook is provided, alerts are gracefully logged to the terminal.

---

## 🚀 How to Run and Test Manually

The project is currently configured with a public test webhook right out of the box! You don't need to configure anything to see the alerts in action.

### Step 1: Run the Application

Open your terminal at the root of the project and simply run:

```bash
go run cmd/gohealthchecker/main.go
```

### Step 2: Observe the Results
- Check your terminal output. You will see the initial health checks running.
- One of the URLs configured in `main.go` (`http://invalid.local.domain`) is guaranteed to fail.
- After **10 seconds** (the second check), the system will trigger a JSON alert payload for the failing service.

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
