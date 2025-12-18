package rules

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

const DefaultMaxConditions = 5

type HighComplexityRule struct {
	MaxConditions int
}

func NewHighComplexityRule() *HighComplexityRule {
	return &HighComplexityRule{
		MaxConditions: DefaultMaxConditions,
	}
}

func (r *HighComplexityRule) ID() string {
	return "GRL006"
}

func (r *HighComplexityRule) Name() string {
	return "high-complexity"
}

func (r *HighComplexityRule) Description() string {
	return "Rules should not have too many conditions in the when clause"
}

func (r *HighComplexityRule) DefaultSeverity() diagnostic.Severity {
	return diagnostic.SeverityWarning
}

func (r *HighComplexityRule) CheckKnowledgeBase(file string, result *parser.ParseResult, kb *ast.KnowledgeBase) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	maxConditions := r.MaxConditions
	if maxConditions == 0 {
		maxConditions = DefaultMaxConditions
	}

	for _, ruleInfo := range result.Rules {
		if ruleInfo.ConditionCount > maxConditions {
			diags = append(diags, diagnostic.Diagnostic{
				File: file,
				Range: diagnostic.Range{
					Start: ruleInfo.WhenPosition,
					End:   ruleInfo.WhenPosition,
				},
				RuleID:   r.ID(),
				RuleName: r.Name(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf("Rule '%s' has %d conditions (max %d) - consider splitting into smaller rules",
					ruleInfo.Name, ruleInfo.ConditionCount, maxConditions),
			})
		}
	}

	return diags
}

var _ SemanticRule = (*HighComplexityRule)(nil)
