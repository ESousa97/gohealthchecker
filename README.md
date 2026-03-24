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

To see the webhook integration in action without needing to create a real Discord/Slack bot right away, you can use a free, temporary webhook service like [Webhook.site](https://webhook.site/).

### Step 1: Get a Test Webhook URL
1. Go to [https://webhook.site/](https://webhook.site/)
2. Copy your **"Your unique URL"** (it looks like `https://webhook.site/your-uuid-here`).

### Step 2: Run the Application

Open your terminal (PowerShell or Bash) at the root of the project and run the following command, replacing the URL with the one you copied:

**Windows (PowerShell):**
```powershell
$env:WEBHOOK_URL="https://webhook.site/YOUR-UUID-HERE"; go run cmd/gohealthchecker/main.go
```

**Linux / macOS:**
```bash
WEBHOOK_URL="https://webhook.site/YOUR-UUID-HERE" go run cmd/gohealthchecker/main.go
```

### Step 3: Observe the Results
- Check your terminal output. You will see initial checks.
- One of the URLs configured in `main.go` (`http://invalid.local.domain`) is guaranteed to fail.
- After **10 seconds** (the second tick), you will see an `[ALERT]` log in the terminal, and if you look at your [Webhook.site](https://webhook.site/) dashboard, you will receive a real JSON payload containing the error!

---

## 🧪 Running Automated Tests

The project includes robust automated tests that utilize Go's `httptest` package to simulate servers and webhook endpoints natively, ensuring everything works without relying on external internet connections.

Run all tests with:
```bash
go test -v ./...
```

You should see an output indicating that the webhook notification logic and the state management (spam prevention) passed successfully.
