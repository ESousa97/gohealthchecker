# GoHealthChecker

> Concurrent and high-performance HTTP health monitoring with a real-time TUI dashboard.

![CI](https://github.com/ESousa97/gohealthchecker/actions/workflows/ci.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gohealthchecker)
![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gohealthchecker.svg)
![License](https://img.shields.io/github/license/ESousa97/gohealthchecker)
![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gohealthchecker)
![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gohealthchecker)

---

`gohealthchecker` is a CLI tool written in Go designed to monitor HTTP endpoints concurrently. It provides a visual dashboard in the terminal (TUI) and sends automatic alerts via Webhooks when services go offline or recover, mitigating false positives through intelligent retries.

## Demonstration

> **TUI Dashboard**: Visualize latency, status, and check times directly in your terminal.

![TUI Screenshot](https://raw.githubusercontent.com/ESousa97/gohealthchecker/main/assets/tui-demo.png)

### Example Usage (As a Library)
```go
c := checker.Checker{
    Targets:  []checker.Target{{URL: "https://google.com"}},
    Notifier: &notifier.ConsoleNotifier{},
}
results := c.Start(5 * time.Second)
```

## Tech Stack

| Technology | Role |
| --- | --- |
| Go 1.25+ | Core language and concurrent execution |
| Bubble Tea | TUI framework for the visual dashboard |
| Viper | Dynamic configuration management (YAML/Env) |
| Lip Gloss | Advanced terminal styling |

## Prerequisites

- Go >= 1.25 (defined in `go.mod`)
- Internet connection for monitoring external URLs

## Installation and Usage

### As a Binary

```bash
go install github.com/ESousa97/gohealthchecker/cmd/gohealthchecker@latest
```

### From Source

```bash
git clone https://github.com/ESousa97/gohealthchecker.git
cd gohealthchecker
cp .env.example .env
# Edit .env (optional for Webhook) or config.yaml
make build
make run
```

## Makefile Targets

| Target | Description |
| --- | --- |
| `make build` | Compiles the binary to `bin/gohealthchecker` |
| `make run` | Runs the application via `go run` |
| `make test` | Runs all unit tests |
| `make lint` | Performs static analysis with golangci-lint |
| `make clean` | Removes build artifacts and temporary files |
| `make install` | Installs the binary to your `GOPATH/bin` |

## Architecture

The project follows principles of extreme modularization and separation of concerns:

- **`internal/checker`**: Core logic for concurrent health checks.
- **`internal/notifier`**: Notification abstraction (Webhook, Console) via interfaces.
- **`internal/tui`**: Visual presentation layer using Bubble Tea.

## API Reference

Detailed documentation of the API and internal packages can be found at [pkg.go.dev](https://pkg.go.dev/github.com/ESousa97/gohealthchecker).

## Configuration

The application uses `config.yaml` or environment variables via Viper.

| Variable | Description | Type | Default |
| --- | --- | --- | --- |
| `WEBHOOK_URL` | Webhook URL for alerts (Slack/Discord) | String | "" |
| `CHECK_INTERVAL`| Interval between checks (e.g., 30s) | Duration | 10s |
| `targets` | List of URLs to monitor (in YAML) | Slice | [] |

## Roadmap

### Completed ✅
- [x] **Phase 1: Synchronized Monitor** - Basic check logic and loop with `time.Ticker`.
- [x] **Phase 2: Concurrency with Channels** - Scaling with goroutines and a central results channel (Fan-out).
- [x] **Phase 3: Notification Layer** - `Notifier` interface and Webhook alerts (Slack/Discord) with anti-spam.
- [x] **Phase 4: Configuration and Resilience** - Viper integration (YAML) and retry logic (3 retries).
- [x] **Phase 5: Observability and TUI** - Interactive visual terminal dashboard with Bubble Tea.

### Next Steps 🚀
- [ ] Support for Telegram and Email notifications
- [ ] Metric export for Prometheus
- [ ] SSL/TLS certificate verification support
- [ ] Persistent uptime history (SQLite)

## Contributing

Contributions are welcome! See the full guide at [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more details.

## Author

**Enoque Sousa**
- Portfolio: [enoquesousa.vercel.app](https://enoquesousa.vercel.app)
- GitHub: [@ESousa97](https://github.com/ESousa97)
