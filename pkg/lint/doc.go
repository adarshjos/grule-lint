// Package lint provides a public API for linting GRL (Grule Rule Language) files.
//
// This package exposes the core linting functionality of grule-lint, allowing
// users to import and use the linter as a library in their own Go applications.
//
// # Basic Usage
//
// Create a linter with default settings and lint a file:
//
//	linter := lint.New()
//	result, err := linter.LintFile("rules.grl")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, d := range result.All() {
//	    fmt.Printf("%s:%d: %s\n", d.File, d.Line(), d.Message)
//	}
//
// # Configuration
//
// Load configuration from a .grl-lint.yaml file:
//
//	cfg, err := lint.LoadConfig(".")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	linter := lint.NewWithConfig(cfg)
//
// # Linting Strings
//
// Lint GRL content directly from a string:
//
//	result := linter.LintString("inline.grl", `
//	    rule MyRule "description" {
//	        when true
//	        then Retract("MyRule");
//	    }
//	`)
//
// # Severity Levels
//
// Diagnostics have four severity levels:
//   - Error: Critical issues that should fail linting
//   - Warning: Potential problems
//   - Info: Informational messages
//   - Hint: Suggestions for improvement
package lint
