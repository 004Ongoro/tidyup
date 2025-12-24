# Contributing to TidyUp 🧹

Thanks for your interest in contributing to **TidyUp**! This document explains how to report issues, propose changes, and get code ready for review so maintainers can merge it quickly.

---

## Table of Contents
- Purpose
- Code of Conduct
- How to report issues
- Feature requests
- Development setup
- Branching & commits
- Pull request process
- Tests, linting & formatting
- Review, CI & merging
- License

---

## Purpose
This file is a short, focused guide to help contributors get started quickly and to set expectations for PRs and issues.

---

## Code of Conduct
We want this to be a welcoming project. Please follow the norms of professional and respectful communication. If you'd like, we can add a `CODE_OF_CONDUCT.md` based on the Contributor Covenant — tell us and we'll add one.

---

## How to report issues
When opening an issue, include:
- A short, descriptive title
- Steps to reproduce the problem (commands run, OS, Go version)
- The command/output or relevant logs
- What you expected vs what happened
- The output of `tidyup --version` and the command you ran (if applicable)

Good issue reports make it faster to triage and fix bugs.

---

## Feature requests
For a new feature, open an issue describing the motivation and a proposed design. If it is a larger change, consider opening a draft PR or discussing the design first in an issue.

---

## Development setup
Minimum requirements:
- Go (reference in `go.mod`) — we recommend Go 1.21+.

Quick start:
```bash
git clone https://github.com/004Ongoro/tidyup.git
cd tidyup
# Run the CLI locally
go run ./cmd scan --path .
# Build a local binary
go build -o tidyup ./cmd
# Run unit tests
go test ./...
```

---

## Branching & commits
- Create branches from `main` using this pattern: `feat/<short-desc>`, `fix/<short-desc>`, `chore/<short-desc>`.
- Use clear commit messages. We recommend the Conventional Commits style:
  - `feat: add new matcher for ...`
  - `fix: avoid scanning node_modules in hidden dirs`
  - `docs: update README`

---

## Pull request process
When opening a PR, please:
- Link to the issue (if any)
- Provide a concise description of the change and why it is needed
- Include test coverage for logic changes where possible
- Run `gofmt -w .` and `go vet ./...` locally

Suggested PR checklist (add to your PR description):
- [ ] I have read the contribution guidelines
- [ ] I have run `go test ./...` and tests pass
- [ ] I have formatted code (`gofmt`) and fixed vet warnings
- [ ] I have added/update tests and documentation where applicable
- [ ] This PR has a clear title and description

---

## Tests, linting & formatting
Please ensure:
- `go test ./...` passes
- Code is formatted with `gofmt -w .`
- Run basic static checks: `go vet ./...` and `staticcheck ./...` (optional)


---

## Review & CI
- Maintainers will review PRs and leave feedback or request changes.
- Small, focused PRs are easier and faster to review.

---

## License
By contributing, you agree that your contributions will be licensed under the repository's MIT license.

---

## Need help?
If you have questions about making a contribution, open an issue and tag a maintainer. Thanks for helping make TidyUp better! 🎉
