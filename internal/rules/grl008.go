package rules

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
)

type UnusedVariableRule struct{}

func (r *UnusedVariableRule) ID() string {
	return "GRL008"
}

func (r *UnusedVariableRule) Name() string {
	return "unused-variable"
}

func (r *UnusedVariableRule) Description() string {
	return "Checks for variables that are assigned but never used"
}

func (r *UnusedVariableRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *UnusedVariableRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, rule := range result.Rules {
		assigned := make(map[string]parser.VariableInfo)
		for _, v := range rule.VariableAssignments {
			assigned[getBaseVarName(v.Name)] = v
		}

		used := make(map[string]bool)
		for _, v := range rule.VariableUsages {
			used[getBaseVarName(v.Name)] = true
		}

		for varName, varInfo := range assigned {
			if !used[varName] {
				diags = append(diags, diagnostic.Diagnostic{
					File:     file,
					Range:    diagnostic.Range{Start: varInfo.Position, End: varInfo.Position},
					RuleID:   r.ID(),
					RuleName: r.Name(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Variable '%s' is assigned but never used in rule '%s'", varName, rule.Name),
				})
			}
		}
	}
	return diags
}
