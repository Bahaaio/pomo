# Contributing to pomo

Thank you for your interest in contributing to pomo!

## Getting Started

1. **Fork the repository** on GitHub

2. **Clone your fork:**

   ```bash
   git clone https://github.com/YOUR_USERNAME/pomo.git
   cd pomo
   ```

3. **Add the upstream remote:**

   ```bash
   git remote add upstream https://github.com/Bahaaio/pomo.git
   ```

4. **Install dependencies:**

   ```bash
   go mod download
   ```

5. **Build the project:**

   ```bash
   go build .
   ```

6. **Run the application:**

   ```bash
   ./pomo
   ```

### Project Structure

```
pomo/
├── .github/         # GitHub Actions workflows
├── actions/         # Post-action command execution
├── cmd/             # CLI commands (Cobra framework)
├── config/          # Configuration loading (Viper)
├── db/              # Database layer (SQLite sessions)
├── ui/              # Terminal UI components (Bubble Tea)
│   ├── ascii/       # ASCII art font rendering
│   ├── colors/      # Color definitions and utilities
│   ├── confirm/     # Confirmation dialog component
│   └── summary/     # Session summary component
└── pomo.go          # Main entry point
```

## Development Workflow

### Creating a Feature Branch

```bash
# Update your local main branch
git switch main
git pull upstream main

# Create a new feature branch
git switch -c feature/your-feature-name
```

### Building and Running

```bash
# Build the project
go build .

# Run directly
go run .

# Run with debug logging
DEBUG=1 go run .
```

### Quality Checks

Before committing, ensure your code passes all checks:

```bash
# Format code
go fmt ./...

# Vet for common issues
go vet ./...

# Run tests
go test ./...

# Tidy dependencies
go mod tidy
```

## Code Style

### Go Conventions

- Follow standard Go conventions and [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (run `go fmt ./...`)
- Run `go vet ./...` to catch common mistakes

### Error Handling

- Always check and handle errors explicitly
- Use descriptive error messages
- Log errors appropriately (use `log.Printf` for debugging)

## Testing

### Writing Tests

- Place tests in `*_test.go` files alongside the code they test
- Use table-driven tests when testing multiple scenarios
- Use descriptive test names that explain what is being tested

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./config
```

### Test Setup

Use `TestMain` to disable logging during tests:

```go
func TestMain(m *testing.M) {
    log.SetOutput(io.Discard)
    os.Exit(m.Run())
}
```

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) for clear, structured commit messages.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring without changing functionality
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, build config, etc.)
- `ci`: CI/CD changes

## Pull Requests

### Before Submitting

1. **Run quality checks:** `go fmt ./...` and `go vet ./...`
2. **Update documentation** if you changed functionality
3. **Rebase on latest main:**

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Creating a Pull Request

1. **Push your branch:**

   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a PR** on GitHub from your fork to `Bahaaio/pomo:main`

3. **Fill out the PR description** with:
   - What changes you made
   - Why you made them
   - How to test them
   - Screenshots (if UI changes)

## Release Process

Releases are automated using GoReleaser when a version tag is pushed.

### Version Tags

- Follow semantic versioning: `v<major>.<minor>.<patch>`
- Current version: `v0.9.0`

## License

By contributing to pomo, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to pomo!
