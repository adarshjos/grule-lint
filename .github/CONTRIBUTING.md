# Contributing to grule-lint

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/grule-lint.git
   cd grule-lint
   ```
3. Install dependencies:
   ```bash
   make deps
   ```

## Development Workflow

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Run All Checks

```bash
make pre-commit
```

## Adding a New Lint Rule

1. Create the rule in `internal/rules/`
2. Register it in `internal/rules/registry.go`
3. Add tests in `test/`
4. Add test fixtures in `testdata/`
5. Update documentation

## Commit Guidelines

- Use clear, descriptive commit messages
- Reference issues when applicable: `Fix #123`
- Keep commits focused and atomic

## Pull Request Process

1. Ensure all tests pass: `make test`
2. Ensure code is formatted: `make fmt`
3. Update documentation if needed
4. Fill out the PR template completely
5. Wait for review

## Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Run `make vet` to check for issues
- Keep functions focused and readable

## Reporting Issues

- Use the issue templates provided
- Include reproduction steps
- Include GRL file samples when relevant

## Questions?

Open an issue for any questions about contributing.
