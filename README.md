<div align="center">
  <h1>gohealthchecker</h1>
  <p>Concurrent and high-performance HTTP health monitoring with a real-time TUI dashboard.</p>

  <img src="assets/tui-demo.png" alt="GoHealthChecker Banner" width="600px">

  <br>

[![CI/CD](https://github.com/ESousa97/gohealthchecker/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/ESousa97/gohealthchecker/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gohealthchecker)](https://goreportcard.com/report/github.com/ESousa97/gohealthchecker)
[![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gohealthchecker.svg)](https://pkg.go.dev/github.com/ESousa97/gohealthchecker)
[![License: MIT](https://img.shields.io/github/license/ESousa97/gohealthchecker)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gohealthchecker)](https://github.com/ESousa97/gohealthchecker)
[![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gohealthchecker)](https://github.com/ESousa97/gohealthchecker/commits/main)

</div>

---

`gohealthchecker` is a CLI tool written in Go designed to monitor HTTP endpoints concurrently. It provides a visual dashboard in the terminal (TUI) and sends automatic alerts via Webhooks when services go offline or recover, mitigating false positives through intelligent retries.

### Example Usage (As a Library)

```go
c := checker.Checker{
    Targets:  []checker.Target{{URL: "https://google.com"}},
    Notifier: &notifier.ConsoleNotifier{},
}
results := c.Start(5 * time.Second)
```

## Tech Stack

| Technology | Role                                        |
| ---------- | ------------------------------------------- |
| Go 1.25+   | Core language and concurrent execution      |
| Bubble Tea | TUI framework for the visual dashboard      |
| Viper      | Dynamic configuration management (YAML/Env) |
| Lip Gloss  | Advanced terminal styling                   |

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

| Target         | Description                                  |
| -------------- | -------------------------------------------- |
| `make build`   | Compiles the binary to `bin/gohealthchecker` |
| `make run`     | Runs the application via `go run`            |
| `make test`    | Runs all unit tests                          |
| `make lint`    | Performs static analysis with golangci-lint  |
| `make clean`   | Removes build artifacts and temporary files  |
| `make install` | Installs the binary to your `GOPATH/bin`     |

## Architecture

The project follows principles of extreme modularization and separation of concerns:

- **`internal/checker`**: Core logic for concurrent health checks.
- **`internal/notifier`**: Notification abstraction (Webhook, Console) via interfaces.
- **`internal/tui`**: Visual presentation layer using Bubble Tea.

## API Reference

Detailed documentation of the API and internal packages can be found at [pkg.go.dev](https://pkg.go.dev/github.com/ESousa97/gohealthchecker).

## Configuration

The application uses `config.yaml` or environment variables via Viper.

| Variable         | Description                            | Type     | Default |
| ---------------- | -------------------------------------- | -------- | ------- |
| `WEBHOOK_URL`    | Webhook URL for alerts (Slack/Discord) | String   | ""      |
| `CHECK_INTERVAL` | Interval between checks (e.g., 30s)    | Duration | 10s     |
| `targets`        | List of URLs to monitor (in YAML)      | Slice    | []      |

## Roadmap

### Completed ✅

- [x] **Phase 1: Synchronized Monitor** - Basic check logic and loop with `time.Ticker`.
- [x] **Phase 2: Concurrency with Channels** - Scaling with goroutines and a central results channel (Fan-out).
- [x] **Phase 3: Notification Layer** - `Notifier` interface and Webhook alerts (Slack/Discord) with anti-spam.
- [x] **Phase 4: Configuration and Resilience** - Viper integration (YAML) and retry logic (3 retries).
- [x] **Phase 5: Observability and TUI** - Interactive visual terminal dashboard with Bubble Tea.

## Contributing

Contributions are welcome! See the full guide at [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more details.

<div align="center">

## Author

**Enoque Sousa**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=flat&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/enoque-sousa-bb89aa168/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=flat&logo=github&logoColor=white)](https://github.com/ESousa97)
[![Portfolio](https://img.shields.io/badge/Portfolio-FF5722?style=flat&logo=target&logoColor=white)](https://enoquesousa.vercel.app)

**[⬆ Back to top](#gohealthchecker)**

Made with ❤️ by [Enoque Sousa](https://github.com/ESousa97)

**Project Status:** Archived — Study Project

</div>
