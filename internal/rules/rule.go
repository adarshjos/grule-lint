// Package rules provides lint rule interfaces and implementations for GRL analysis.
package rules

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

// Rule is the base interface for all lint rules.
type Rule interface {
	// ID returns the unique identifier for this rule (e.g., "GRL001").
	ID() string

	// Name returns the human-readable name (e.g., "syntax-error").
	Name() string

	// Description returns a brief description of what this rule checks.
	Description() string

	// DefaultSeverity returns the default severity level for this rule.
	DefaultSeverity() diagnostic.Severity
}

// SyntaxRule is a rule that runs when parsing fails (no AST available).
// These rules analyze parse errors.
type SyntaxRule interface {
	Rule

	// CheckSyntax analyzes parse errors and returns diagnostics.
	CheckSyntax(result *parser.ParseResult) []diagnostic.Diagnostic
}

// SemanticRule is a rule that runs when parsing succeeds (AST available).
// These rules analyze the AST for semantic issues.
type SemanticRule interface {
	Rule

	// CheckKnowledgeBase analyzes the parsed AST and returns diagnostics.
	// It receives the full ParseResult which includes accurate position info.
	CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic
}
