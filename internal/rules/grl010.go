package rules

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
)

type EmptyWhenRule struct{}

func (r *EmptyWhenRule) ID() string {
	return "GRL010"
}

func (r *EmptyWhenRule) Name() string {
	return "empty-when"
}

func (r *EmptyWhenRule) Description() string {
	return "Checks for rules with empty when clauses (always true)"
}

func (r *EmptyWhenRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *EmptyWhenRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, rule := range result.Rules {
		// Check for empty when or always-true conditions:
		// 1. No when expression at all
		// 2. When expression is just 'true' literal (no actual comparisons)
		isEmpty := !rule.HasWhenExpression && rule.ConditionCount == 0
		isAlwaysTrue := rule.ConditionCount == 0 && (rule.WhenExpressionText == "true" || rule.WhenExpressionText == "TRUE")

		if isEmpty || isAlwaysTrue {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: rule.WhenPosition,
					End:   rule.WhenPosition,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule '%s' has an empty or always-true when clause - it will always fire", rule.Name),
			})
		}
	}

	return diags
}
