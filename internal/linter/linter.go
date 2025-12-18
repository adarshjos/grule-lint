// Package linter provides the main linting orchestration for GRL files.
package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	"github.com/adarshjos/grule-lint/internal/rules"
)

// Linter orchestrates the linting process.
type Linter struct {
	parser   *parser.Parser
	registry *rules.Registry
}

// New creates a new Linter with the default registry.
func New() *Linter {
	return &Linter{
		parser:   parser.NewParser(),
		registry: rules.DefaultRegistry(),
	}
}

// NewWithRegistry creates a new Linter with a custom registry.
func NewWithRegistry(registry *rules.Registry) *Linter {
	return &Linter{
		parser:   parser.NewParser(),
		registry: registry,
	}
}

// NewWithConfig creates a new Linter with the specified configuration.
func NewWithConfig(cfg rules.RegistryConfig) *Linter {
	return &Linter{
		parser:   parser.NewParser(),
		registry: rules.DefaultRegistryWithConfig(cfg),
	}
}

// LintFile lints a single GRL file.
func (l *Linter) LintFile(file string) (*diagnostic.DiagnosticSet, error) {
	result, err := l.parser.ParseFile(file)
	if err != nil {
		return nil, err
	}

	return l.lintParseResult(result), nil
}

// LintString lints GRL content from a string.
func (l *Linter) LintString(file, content string) *diagnostic.DiagnosticSet {
	result := l.parser.ParseString(file, content)
	return l.lintParseResult(result)
}

// lintParseResult runs all applicable rules on a parse result.
func (l *Linter) lintParseResult(result *parser.ParseResult) *diagnostic.DiagnosticSet {
	ds := diagnostic.NewDiagnosticSet()

	if len(result.Errors) > 0 {
		// Run syntax rules on parse errors
		diags := l.registry.RunSyntaxRules(result)
		ds.AddAll(diags)
	} else {
		// No parse errors - run semantic rules
		// Semantic rules use result.Rules (from ANTLR) not necessarily the KB
		diags := l.registry.RunSemanticRules(result)
		ds.AddAll(diags)
	}

	return ds
}

// LintFiles lints multiple GRL files.
func (l *Linter) LintFiles(files []string) (*diagnostic.DiagnosticSet, error) {
	ds := diagnostic.NewDiagnosticSet()

	for _, file := range files {
		fileDiags, err := l.LintFile(file)
		if err != nil {
			return nil, err
		}
		ds.AddAll(fileDiags.All())
	}

	return ds, nil
}

// LintDirectory lints all GRL files in a directory (recursively).
func (l *Linter) LintDirectory(dir string) (*diagnostic.DiagnosticSet, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking directory: %w", err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".grl") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scanning directory %s: %w", dir, err)
	}

	return l.LintFiles(files)
}

// LintPaths lints files and/or directories.
func (l *Linter) LintPaths(paths []string) (*diagnostic.DiagnosticSet, error) {
	ds := diagnostic.NewDiagnosticSet()

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("accessing path %s: %w", path, err)
		}

		var pathDiags *diagnostic.DiagnosticSet
		if info.IsDir() {
			pathDiags, err = l.LintDirectory(path)
		} else {
			pathDiags, err = l.LintFile(path)
		}

		if err != nil {
			return nil, err
		}
		ds.AddAll(pathDiags.All())
	}

	return ds, nil
}
