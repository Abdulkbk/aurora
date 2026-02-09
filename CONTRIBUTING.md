# Contributing to Aurora 🌌

Thank you for your interest in contributing to Aurora! This document
outlines the guidelines and Go best practices to follow when submitting
changes.

## Getting Started

1. Fork the repository and clone your fork.
2. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```
3. Make your changes, following the guidelines below.
4. Run checks before committing:
   ```bash
   make fmt
   make test
   ```
5. Commit with a clear message and open a Pull Request.

## Go Best Practices

### Code Style

- **Follow `gofmt`** — All code must be formatted with `go fmt`. Run
  `make fmt` before committing.
- **Line length** — Keep lines to **80 characters** where reasonable.
  The project enforces an 80-column ruler in `.vscode/settings.json`.
- **Tab size** — Use **tabs** with a width of 8 for indentation (Go
  standard).

### Naming Conventions

- Use **MixedCaps** (`ExportedName`) and **mixedCaps** (`unexportedName`).
  Do not use underscores in Go names.
- Keep names **short and descriptive**. Prefer `pr` over
  `pullRequest` for local variables, but use clear names for exported
  symbols.
- Acronyms should be **all caps**: `URL`, `PR`, `HTTP`, `ID` — not
  `Url`, `Pr`, `Http`, `Id`.
- Interface names should describe behavior, e.g. `Builder`, `Parser`.
  Do not prefix with `I`.

### Error Handling

- **Always check errors.** Never discard an error with `_`.
- **Wrap errors with context** using `fmt.Errorf`:
  ```go
  if err != nil {
      return fmt.Errorf("failed to fetch PR details: %w", err)
  }
  ```
- Return errors to the caller rather than logging and continuing.
  Let the top-level command handle user-facing output.

### Functions & Methods

- Keep functions **short and focused** — each function should do one
  thing.
- Prefer **returning errors** over `log.Fatal` or `os.Exit` in library
  code. Only the `main` package or CLI command handlers should exit.
- Use **named return values** sparingly — only when they improve
  readability.

### Packages & Project Structure

Aurora follows the standard Go project layout:

```
cmd/            # CLI commands (cobra)
internal/       # Private application code
  docker/       # Docker image building logic
  github/       # GitHub API client
main.go         # Entry point
```

- Code in `internal/` is **not importable** by external packages —
  keep implementation details here.
- Each package should have a **clear, single responsibility**.
- Avoid circular dependencies between packages.

### Comments & Documentation

- All **exported** functions, types, and constants must have a doc
  comment starting with the name:
  ```go
  // ParsePRURL extracts the owner, repo, and PR number from a
  // GitHub pull request URL.
  func ParsePRURL(rawURL string) (*PRInfo, error) {
  ```
- Use `//` comments, not `/* */`, for regular code comments.
- Write comments that explain **why**, not what.

### Testing

- Place tests in `_test.go` files in the same package.
- Use table-driven tests where there are multiple cases:
  ```go
  tests := []struct {
      name    string
      input   string
      want    string
      wantErr bool
  }{
      {"valid PR URL", "https://...", "expected", false},
      {"invalid URL", "not-a-url", "", true},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          // ...
      })
  }
  ```
- Run tests with `make test` before submitting a PR.
- Aim for meaningful tests — test behavior, not implementation.

### Dependencies

- Keep dependencies **minimal**. Aurora currently only depends on
  `cobra` for CLI.
- Run `make tidy` to clean up `go.mod` and `go.sum` after adding or
  removing imports.
- Avoid adding large frameworks when the standard library suffices.

### Concurrency

- Prefer channels and `sync` primitives over shared mutable state.
- When using goroutines, always ensure they are properly cleaned up
  (use `context.Context` for cancellation).
- Protect shared state with `sync.Mutex` — document what the mutex
  guards.

## Commit Messages

Write clear, descriptive commit messages:

```
<area>: <short summary in lowercase>

Optional longer description explaining the motivation
and any non-obvious design decisions.
```

Examples:

```
docker: add support for Core Lightning builds
github: handle rate limiting in PR fetcher
cmd: validate --tag flag format
```

## Code Review

All submissions require review before merging. Reviewers will look
for:

- Adherence to these guidelines
- Test coverage for new functionality
- Clear error messages and user-facing output
- No regressions in existing functionality

## Questions?

If you're unsure about anything, open an issue to discuss your idea
before writing code. We're happy to help!
