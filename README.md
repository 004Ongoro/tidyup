# TidyUp 🧹

**TidyUp** is a fast, safety-first CLI written in Go that helps developers reclaim disk space by locating and removing stale dependency or build directories (e.g., `node_modules`, `target`, `.venv`). It favors conservative defaults and explicit user consent so you can automate cleanup with confidence.

[![Go Version](https://img.shields.io/github/go-mod/go-version/004Ongoro/tidyup)](https://go.dev/)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Table of Contents

- [Key Features](#-key-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [CLI Reference](#-cli-reference)
- [Safety & Design](#-safety--design)
- [Scheduling Behavior](#-scheduling-behavior)
- [Development & Contributing](#-development--contributing)
- [Troubleshooting](#-troubleshooting)
- [License](#-license)

---

## 🚀 Key Features

- **Parallel Scanning**: Worker-pool scanning across CPU cores for speed.
- **Deep Scan Mode**: Optionally ignore anchor-file checks and rely purely on directory metadata.
- **Interactive Cleanup**: Multi-select TUI to confirm deletion of targets.
- **Automated Scheduling**: Native Task Scheduler (Windows) or Cron (macOS/Linux).
- **YAML Configuration**: Customize matchers and blocklists via `.tidyup.yaml`.
- **Safety-First Defaults**: Built-in protections for system and user-critical folders.

---

## 📦 Installation

Install the latest released binary via go:

```bash
go install github.com/004Ongoro/tidyup@latest
```

Build from source:

```bash
git clone https://github.com/004Ongoro/tidyup.git
cd tidyup
go build ./...
# Or build a single binary
go build -o tidyup cmd/main.go
```

Run locally without installing:

```bash
go run ./cmd -h
```

---

## ⚡ Quick Start

1. Preview stale folders (default: 30 days):

```bash
tidyup scan --path "C:/Users/You/Projects" --days 30
```

2. Interactively delete selected folders:

```bash
tidyup clean --path . --days 30
```

3. Force delete without prompts (use with caution):

```bash
tidyup clean --path . --days 30 --force
```

4. Set up automated cleaning (daily or weekly):

```bash
tidyup schedule --interval daily
# Remove the schedule
tidyup schedule --remove
```

---

## Safe Cleanup(For Testing)

```bash
# Simulate a cleanup without deleting files
tidyup clean --dry-run

# Force a cleanup (useful for automated tasks)
tidyup clean --force
```

## ⚙️ Configuration

Create a `.tidyup.yaml` at `$HOME` or project root to override defaults.

Example `.tidyup.yaml`:

```yaml
blocklist:
  - "DoNotTouch"
  - "Backups"

matchers:
  - name: "Node.js"
    target_dir: "node_modules"
    anchor_file: "package.json"
  - name: "Rust"
    target_dir: "target"
    anchor_file: "Cargo.toml"
```

- `blocklist` — directories that will never be considered for deletion.
- `matchers` — set `target_dir` that identifies the directory name to match and an `anchor_file` used to determine project activity by inspecting its modification time.

You can also pass a custom config file with `--config /path/to/.tidyup.yaml`.

---


## 📋 Testing

TidyUp includes a robust test suite to ensure path safety across different operating systems.

```bash

go test ./cmd -v
```

## 🔧 CLI Reference

Commands (summary):

- `tidyup scan` — Scan for stale project directories and print a summary.
  - `--path, -p` (default `.`) — root directory to scan
  - `--days, -d` (default `30`) — age threshold in days
  - `--deep` — ignore anchor files and use directory metadata only

- `tidyup clean` — Interactively remove discovered directories.
  - `--path, -p` — root directory to scan
  - `--days, -d` — age threshold
  - `--force, -f` — skip confirmation and delete all found targets
  - `--deep` — deep mode

- `tidyup schedule` — Create or remove a scheduled background task.
  - `--interval, -i` (`daily` or `weekly`)
  - `--remove, -r` — remove scheduled task

- `--config` — provide an alternate config file
- `--version` — print version

For more information about flags, run `tidyup <command> --help`.

---

## 🛡 Safety & Design

- TidyUp is **conservative by default**: it will only propose directories that match configured matchers and pass age/anchor checks.
- It **skips system and common OS folders** to prevent accidental data loss (e.g., Program Files, AppData, /System on macOS).
- Use `--dry-run` patterns (via `scan`) and `--force` only when you are sure.
- Add sensitive folders to `blocklist` in `.tidyup.yaml` to guarantee they are never touched.

---

## 🗓 Scheduling Behavior

- Windows: Uses `schtasks` to create a `TidyUpAutoClean` task that runs `tidyup clean --force` on the requested schedule.
- macOS/Linux: Appends a cron entry to run `tidyup clean --force` at noon daily (or weekly on Sundays).

Note: On Unix, the schedule is appended to the user's crontab. On Windows, the task is created in Task Scheduler; administrative privileges may be required.

---

## 👩‍💻 Development & Contributing

Thanks for your interest in contributing! Please read the following before opening issues or PRs.

Prerequisites:

- Go 1.21+ (or the version in `go.mod`)
- `gofmt` (Go tools)

Local development workflow:

```bash
git clone https://github.com/004Ongoro/tidyup.git
cd tidyup
# Run the CLI locally
go run ./cmd scan --path .
# Build a local binary
go build -o tidyup ./cmd
```

Testing & Quality:

- Run tests (if any): `go test ./...`
- Format code: `gofmt -w .`
- Static checks: `go vet ./...` and `staticcheck ./...` (if installed)

Contributing guidelines:

1. Open an issue for any bug or feature request with a clear description and reproduction steps.
2. Create a branch from `main` using the pattern `feat/<short-description>` or `fix/<short-description>`.
3. Run tests and code formatters locally before opening a PR.
4. Include details in your PR about why the change is needed and any migration notes.

Maintainers will review and request changes if needed.

---

## ❓ Troubleshooting

- Permission errors when scheduling on Windows: run the terminal as Administrator.
- Cron job not running: ensure `tidyup` is in your PATH or use an absolute path to the binary in the cron entry.
- False positives: add non-target folders to `blocklist` and adjust `matchers` in `.tidyup.yaml`.

If you still need help, open an issue and include `tidyup --version` and the command you ran.

---

## 📄 License

Distributed under the **MIT License**. See `LICENSE` for details.

