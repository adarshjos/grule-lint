package rules

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
)

type EmptyThenRule struct{}

func (r *EmptyThenRule) ID() string {
	return "GRL011"
}

func (r *EmptyThenRule) Name() string {
	return "empty-then"
}

func (r *EmptyThenRule) Description() string {
	return "Checks for rules with empty then clauses (no actions)"
}

func (r *EmptyThenRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *EmptyThenRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, rule := range result.Rules {
		if len(rule.FunctionCalls) == 0 && rule.ThenActionCount == 0 {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: rule.ThenPosition,
					End:   rule.ThenPosition,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule '%s' has an empty then clause - it performs no actions", rule.Name),
			})
		}
	}

	return diags
}
