package rules

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/iancoleman/strcase"
)

type NamingConvention string

const (
	ConventionPascalCase NamingConvention = "PascalCase"
	ConventionCamelCase  NamingConvention = "camelCase"
	ConventionSnakeCase  NamingConvention = "snake_case"
	ConventionKebabCase  NamingConvention = "kebab-case"
)

type NamingConventionRule struct {
	Convention NamingConvention
}

func NewNamingConventionRule() *NamingConventionRule {
	return &NamingConventionRule{
		Convention: ConventionPascalCase,
	}
}

// NewNamingConventionRuleWithConfig creates a NamingConventionRule with the specified convention.
func NewNamingConventionRuleWithConfig(convention string) *NamingConventionRule {
	r := &NamingConventionRule{
		Convention: ConventionPascalCase,
	}
	r.SetConvention(convention)
	return r
}

// SetConvention sets the naming convention to enforce.
// Valid values: "PascalCase", "camelCase", "snake_case", "kebab-case"
func (r *NamingConventionRule) SetConvention(convention string) {
	switch convention {
	case "PascalCase":
		r.Convention = ConventionPascalCase
	case "camelCase":
		r.Convention = ConventionCamelCase
	case "snake_case":
		r.Convention = ConventionSnakeCase
	case "kebab-case":
		r.Convention = ConventionKebabCase
	}
}

func (r *NamingConventionRule) ID() string {
	return "GRL007"
}

func (r *NamingConventionRule) Name() string {
	return "naming-convention"
}

func (r *NamingConventionRule) Description() string {
	return "Checks that rule names follow a consistent naming convention"
}

func (r *NamingConventionRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *NamingConventionRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *gruleAst.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, rule := range result.Rules {
		if !r.isValidName(rule.Name) {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: rule.Position,
					End:   rule.Position,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Rule name '%s' does not follow %s convention", rule.Name, r.Convention),
			})
		}
	}
	return diags
}

func (r *NamingConventionRule) isValidName(name string) bool {
	if name == "" {
		return false
	}
	switch r.Convention {
	case ConventionPascalCase:
		return name == strcase.ToCamel(name)
	case ConventionCamelCase:
		return name == strcase.ToLowerCamel(name)
	case ConventionSnakeCase:
		return name == strcase.ToSnake(name)
	case ConventionKebabCase:
		return name == strcase.ToKebab(name)
	default:
		return true
	}
}
