package rules

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
)

type ConflictingRulesRule struct{}

func (r *ConflictingRulesRule) ID() string {
	return "GRL012"
}

func (r *ConflictingRulesRule) Name() string {
	return "conflicting-rules"
}

func (r *ConflictingRulesRule) Description() string {
	return "Checks for rules with identical conditions but different actions"
}

func (r *ConflictingRulesRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *ConflictingRulesRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	whenGroups := make(map[string][]parser.RuleInfo)
	for _, rule := range result.Rules {
		if rule.WhenExpressionText != "" {
			whenGroups[rule.WhenExpressionText] = append(whenGroups[rule.WhenExpressionText], rule)
		}
	}

	for _, rules := range whenGroups {
		if len(rules) < 2 {
			continue
		}

		salienceSet := make(map[string]bool)
		for _, rule := range rules {
			salienceSet[rule.Salience] = true
		}

		if len(salienceSet) == 1 || (len(salienceSet) == 0) {
			for i := 1; i < len(rules); i++ {
				diags = append(diags, diagnostic.Diagnostic{
					File: file,
					Range: diagnostic.Range{
						Start: rules[i].Position,
						End:   rules[i].Position,
					},
					RuleID:   r.ID(),
					RuleName: r.Name(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Rule '%s' has the same when clause as '%s' (at line %d) - this may cause conflicts",
						rules[i].Name,
						rules[0].Name,
						rules[0].Position.Line,
					),
				})
			}
		} else {
			for i := 1; i < len(rules); i++ {
				diags = append(diags, diagnostic.Diagnostic{
					File: file,
					Range: diagnostic.Range{
						Start: rules[i].Position,
						End:   rules[i].Position,
					},
					RuleID:   r.ID(),
					RuleName: r.Name(),
					Severity: diagnostic.SeverityInfo,
					Message: fmt.Sprintf(
						"Rule '%s' has the same when clause as '%s' (at line %d) - salience differs so execution order is defined",
						rules[i].Name,
						rules[0].Name,
						rules[0].Position.Line,
					),
				})
			}
		}
	}

	return diags
}
