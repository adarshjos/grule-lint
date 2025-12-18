# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - TBD

### Added
- Initial release of grule-lint
- 12 lint rules for GRL files:
  - GRL001: syntax-error - Detects invalid GRL syntax
  - GRL002: missing-description - Warns when rules lack descriptions
  - GRL003: missing-salience - Info when rules have no explicit salience
  - GRL004: missing-retract - Warns about potential infinite loops
  - GRL005: duplicate-rule - Errors on duplicate rule names
  - GRL006: high-complexity - Warns when conditions exceed threshold
  - GRL007: naming-convention - Enforces PascalCase naming
  - GRL008: unused-variable - Detects unused variable assignments
  - GRL009: undefined-variable - Hints about undefined variables
  - GRL010: empty-when - Warns on always-true conditions
  - GRL011: empty-then - Warns on rules with no actions
  - GRL012: conflicting-rules - Detects rules with identical conditions
- CLI with configurable options:
  - `--config` for custom configuration file
  - `--exclude` for file pattern exclusion
  - `--quiet` for error-only output
  - `--no-color` for plain text output
  - `--output` for file output
- YAML configuration support via `.grl-lint.yaml`
- Text output format with colored severity indicators
- GitHub Actions CI/CD workflows
- GoReleaser configuration for cross-platform releases

[Unreleased]: https://github.com/adarshjos/grule-lint/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/adarshjos/grule-lint/releases/tag/v0.1.0
