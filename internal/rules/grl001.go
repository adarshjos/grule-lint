package rules

import (
	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

type SyntaxErrorRule struct{}

func (r *SyntaxErrorRule) ID() string {
	return "GRL001"
}

func (r *SyntaxErrorRule) Name() string {
	return "syntax-error"
}

func (r *SyntaxErrorRule) Description() string {
	return "Reports syntax errors in GRL files that prevent parsing"
}

func (r *SyntaxErrorRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityError
}

// CheckSyntax converts parse errors to diagnostics.
func (r *SyntaxErrorRule) CheckSyntax(result *parser.ParseResult) []diagnostic.Diagnostic {
	if result.Success() {
		return nil
	}

	var diags []diagnostic.Diagnostic

	for _, err := range result.Errors {
		diags = append(diags, diagnostic.Diagnostic{
			File: result.File,
			Range: diagnostic.Range{
				Start: diagnostic.Position{Line: err.Line, Column: err.Column},
				End:   diagnostic.Position{Line: err.Line, Column: err.Column},
			},
			RuleID:   r.ID(),
			RuleName: r.Name(),
			Severity: r.DefaultSeverity(),
			Message:  err.Message,
		})
	}

	return diags
}

var _ SyntaxRule = (*SyntaxErrorRule)(nil)
