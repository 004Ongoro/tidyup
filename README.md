# TidyUp 🧹

**TidyUp** is a blazing-fast, safety-first CLI tool built in Go to help developers reclaim gigabytes of disk space. It identifies and removes "stale" dependency folders (like `node_modules`, `target`, and `.venv`) based on the last time you actually worked on a project.

[![Go Version](https://img.shields.io/github/go-mod/go-version/004Ongoro/tidyup)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 Features

* **Eco-system Aware**: Automatically detects Node.js (`node_modules`), Rust (`target`), Python (`venv`), Maven/Gradle (`build`), and more.
* **Smart Age Detection**: Instead of just looking at the folder age, TidyUp checks "Anchor Files" (like `package.json`) to see when you last modified the project.
* **Safety First**: Built-in blocklist prevents accidental deletion of system files, editor extensions (VS Code), or critical AppData.
* **Interactive Mode**: Use a checkbox-style interface to choose exactly what to delete.
* **High Performance**: Leverages Go's concurrency (goroutines) to scan your drive in seconds.
* **Visual Feedback**: Colorful and intuitive terminal output to easily distinguish between project types and sizes.

## 📦 Installation

Ensure you have [Go](https://go.dev/dl/) installed, then run:

```bash
go install [github.com/004Ongoro/tidyup@latest](https://github.com/004Ongoro/tidyup@latest)

```

Alternatively, clone and build locally:

```bash
git clone [https://github.com/004Ongoro/tidyup.git](https://github.com/004Ongoro/tidyup.git)
cd tidyup
go build -o tidyup

```

## 🛠 Usage

### 1. Scan for stale projects

Preview how much space you can save without deleting anything.

```bash
# Scan current directory for projects untouched for 30 days
tidyup scan

# Scan a specific path for projects untouched for 60 days
tidyup scan --path C:/Users/Name/Projects --days 60

```

### 2. Clean up space

Run an interactive cleanup where you select which folders to remove via a multi-select menu.

```bash
tidyup clean --path . --days 30

```

### 3. Force cleanup

Skip the interactive menu and delete everything matching the criteria (use with caution!).

```bash
tidyup clean --path . --days 90 --force

```

## 🛡 Safety & Blocklist

TidyUp is designed to never touch your operating system or installed applications. It automatically ignores:

* **Windows**: `AppData`, `Program Files`, `Windows`, `System32`
* **macOS/Linux**: `Library`, `.cache`
* **Editors**: `.vscode`, `.antigravity`, `.cursor`
* **Toolchains**: `.rustup`, `.cargo`

## 📋 Supported Ecosystems

| Language | Target Folder | Anchor File |
| --- | --- | --- |
| **Node.js** | `node_modules` | `package.json` |
| **Rust** | `target` | `Cargo.toml` |
| **Python** | `.venv`, `venv` | `pyproject.toml`, `requirements.txt` |
| **Java (Maven)** | `target` | `pom.xml` |
| **Java (Gradle)** | `build` | `build.gradle` |

## 🗺 Roadmap

* [ ] **Configuration File**: Support for `.tidyup.yaml` to save custom blocklists and paths.
* [ ] **Custom Matchers**: Allow users to define their own target/anchor pairs via the config.
* [ ] **Deep Scan**: Option to ignore `.gitignore` rules for a more thorough cleanup.
* [ ] **Automated Scheduling**: Support for running as a background cron job or scheduled task.

## 🤝 Contributing

Contributions are welcome! Whether it's adding a new project matcher, improving the UI, or fixing a bug:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

```

```