package rules

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

type MissingSalienceRule struct{}

func (r *MissingSalienceRule) ID() string {
	return "GRL003"
}

func (r *MissingSalienceRule) Name() string {
	return "missing-salience"
}

func (r *MissingSalienceRule) Description() string {
	return "Rules should have an explicit salience value to clarify execution order"
}

func (r *MissingSalienceRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityInfo
}

func (r *MissingSalienceRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, ruleInfo := range result.Rules {
		if ruleInfo.Salience == "" {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: ruleInfo.Position,
					End:   ruleInfo.Position,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule '%s' does not specify salience (defaults to 0)", ruleInfo.Name),
			})
		}
	}

	return diags
}

var _ SemanticRule = (*MissingSalienceRule)(nil)
