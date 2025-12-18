package rules

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

type MissingDescriptionRule struct{}

func (r *MissingDescriptionRule) ID() string {
	return "GRL002"
}

func (r *MissingDescriptionRule) Name() string {
	return "missing-description"
}

func (r *MissingDescriptionRule) Description() string {
	return "Rules should have a description explaining their purpose"
}

func (r *MissingDescriptionRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *MissingDescriptionRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, ruleInfo := range result.Rules {
		if ruleInfo.Description == "" {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: ruleInfo.Position,
					End:   ruleInfo.Position,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule '%s' is missing a description", ruleInfo.Name),
			})
		}
	}

	return diags
}

var _ SemanticRule = (*MissingDescriptionRule)(nil)
