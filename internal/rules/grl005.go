package rules

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

type DuplicateRuleRule struct{}

func (r *DuplicateRuleRule) ID() string {
	return "GRL005"
}

func (r *DuplicateRuleRule) Name() string {
	return "duplicate-rule"
}

func (r *DuplicateRuleRule) Description() string {
	return "Rule names must be unique within a file"
}

func (r *DuplicateRuleRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityError
}

func (r *DuplicateRuleRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	seen := make(map[string]diagnostic.Position)

	for _, ruleInfo := range result.Rules {
		if firstPos, exists := seen[ruleInfo.Name]; exists {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: ruleInfo.Position,
					End:   ruleInfo.Position,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Duplicate rule name '%s' (first defined at line %d)", ruleInfo.Name, firstPos.Line),
			})
		} else {
			seen[ruleInfo.Name] = ruleInfo.Position
		}
	}

	return diags
}

var _ SemanticRule = (*DuplicateRuleRule)(nil)
