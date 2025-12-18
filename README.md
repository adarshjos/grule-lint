# grule-lint

[![Go Version](https://img.shields.io/github/go-mod/go-version/adarshjos/grule-lint)](https://go.dev/)
[![CI](https://github.com/adarshjos/grule-lint/actions/workflows/ci.yml/badge.svg)](https://github.com/adarshjos/grule-lint/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/adarshjos/grule-lint)](https://github.com/adarshjos/grule-lint/releases)
[![License](https://img.shields.io/github/license/adarshjos/grule-lint)](LICENSE)

A static analysis tool for GRL (Grule Rule Language) files used with the [Grule Rule Engine](https://github.com/hyperjumptech/grule-rule-engine).

## Features

- Detects syntax errors before runtime
- Enforces best practices for rule definitions
- Configurable lint rules
-  output formats (text)
## Lint Rules

| Rule ID | Name | Description |
|---------|------|-------------|
| GRL001 | syntax-error | Invalid GRL syntax |
| GRL002 | missing-description | Rule lacks a description |
| GRL003 | missing-salience | Rule has no explicit salience |
| GRL004 | missing-retract | Rule doesn't call Retract() (potential infinite loop) |
| GRL005 | duplicate-rule | Duplicate rule name detected |
| GRL006 | high-complexity | When clause has too many conditions |
| GRL007 | naming-convention | Rule name doesn't follow convention |
| GRL008 | empty-when | When clause is empty |
| GRL009 | conflicting-rules | Rules with same conditions but different actions |

## Installation

### Go Install

```bash
go install github.com/adarshjos/grule-lint/cmd/grule-lint@latest
```

### From Source

```bash
git clone https://github.com/adarshjos/grule-lint.git
cd grule-lint
make build
```

### Download Binary

Download pre-built binaries from the [Releases](https://github.com/adarshjos/grule-lint/releases) page.

## Usage

```bash
# Lint a single file
grule-lint rules.grl

# Lint a directory
grule-lint rules/

# Lint with glob pattern
grule-lint rules/*.grl

# Use custom config
grule-lint --config .grl-lint.yaml rules/

# Output as JSON
grule-lint --output results.json rules/

# Exclude patterns
grule-lint --exclude "**/test/**" rules/
```

## Configuration

Create a `.grl-lint.yaml` file in your project root:

```yaml
rules:
  GRL001: error      # syntax-error (cannot be disabled)
  GRL002: warning    # missing-description
  GRL003: warning    # missing-salience
  GRL004: error      # missing-retract
  GRL005: error      # duplicate-rule
  GRL006: warning    # high-complexity
  GRL007: warning    # naming-convention
  GRL008: warning    # empty-when
  GRL009: warning    # conflicting-rules

settings:
  naming-convention:
    pattern: "^[A-Z][a-zA-Z0-9_]*$"
  high-complexity:
    max-conditions: 5

exclude:
  - "**/vendor/**"
  - "**/testdata/**"
```

## Output Formats

### Text (default)
```
rules/order.grl:15:1: GRL004 [warning] Rule 'ProcessOrder' does not call Retract()
```

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
