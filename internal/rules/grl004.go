package rules

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

type MissingRetractRule struct{}

func (r *MissingRetractRule) ID() string {
	return "GRL004"
}

func (r *MissingRetractRule) Name() string {
	return "missing-retract"
}

func (r *MissingRetractRule) Description() string {
	return "Rules should call Retract() to prevent infinite loops"
}

func (r *MissingRetractRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *MissingRetractRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, ruleInfo := range result.Rules {
		hasRetract := false

		for _, fc := range ruleInfo.FunctionCalls {
			if fc.Name == "Retract" {
				hasRetract = true
				break
			}
		}

		if !hasRetract {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: ruleInfo.ThenPosition,
					End:   ruleInfo.ThenPosition,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule '%s' does not call Retract() - this may cause an infinite loop", ruleInfo.Name),
			})
		}
	}

	return diags
}

var _ SemanticRule = (*MissingRetractRule)(nil)
