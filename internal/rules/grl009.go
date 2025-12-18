package rules

import (
	"fmt"
	"strings"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
)

type UndefinedVariableRule struct {
	KnownDataContextNames []string
}

func NewUndefinedVariableRule() *UndefinedVariableRule {
	return &UndefinedVariableRule{
		KnownDataContextNames: []string{
			"Order", "Customer", "Product", "Cart", "User", "Request",
			"Response", "Input", "Output", "Data", "Context", "Config",
			"Fact", "Facts", "Result", "Results", "Item", "Items",
		},
	}
}

func (r *UndefinedVariableRule) ID() string {
	return "GRL009"
}

func (r *UndefinedVariableRule) Name() string {
	return "undefined-variable"
}

func (r *UndefinedVariableRule) Description() string {
	return "Checks for variables that are used but may not be defined"
}

func (r *UndefinedVariableRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityHint
}

func (r *UndefinedVariableRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	knownNames := make(map[string]bool)
	for _, name := range r.KnownDataContextNames {
		knownNames[strings.ToLower(name)] = true
	}

	for _, rule := range result.Rules {
		localVars := make(map[string]bool)
		for _, v := range rule.VariableAssignments {
			localVars[strings.ToLower(getBaseVarName(v.Name))] = true
		}

		seen := make(map[string]bool)
		for _, v := range rule.VariableUsages {
			baseName := getBaseVarName(v.Name)
			baseNameLower := strings.ToLower(baseName)

			if seen[baseNameLower] || knownNames[baseNameLower] || localVars[baseNameLower] || builtIns[baseName] {
				continue
			}

			seen[baseNameLower] = true
			diags = append(diags, diagnostic.Diagnostic{
				File:     file,
				Range:    diagnostic.Range{Start: v.Position, End: v.Position},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Variable '%s' may not be defined in rule '%s'", baseName, rule.Name),
			})
		}
	}
	return diags
}

var builtIns = map[string]bool{
	"Retract": true, "Log": true, "Now": true, "IsNil": true, "IsZero": true,
	"Len": true, "MakeTime": true, "Changed": true, "Complete": true,
	"true": true, "false": true, "nil": true,
}
