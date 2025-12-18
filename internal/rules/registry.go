package rules

import (
	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/parser"
)

// Registry holds all registered lint rules.
type Registry struct {
	syntaxRules   []SyntaxRule
	semanticRules []SemanticRule
	allRules      map[string]Rule
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		syntaxRules:   make([]SyntaxRule, 0),
		semanticRules: make([]SemanticRule, 0),
		allRules:      make(map[string]Rule),
	}
}

// RegisterSyntax registers a syntax rule.
func (r *Registry) RegisterSyntax(rule SyntaxRule) {
	r.syntaxRules = append(r.syntaxRules, rule)
	r.allRules[rule.ID()] = rule
}

// RegisterSemantic registers a semantic rule.
func (r *Registry) RegisterSemantic(rule SemanticRule) {
	r.semanticRules = append(r.semanticRules, rule)
	r.allRules[rule.ID()] = rule
}

// SyntaxRules returns all registered syntax rules.
func (r *Registry) SyntaxRules() []SyntaxRule {
	return r.syntaxRules
}

// SemanticRules returns all registered semantic rules.
func (r *Registry) SemanticRules() []SemanticRule {
	return r.semanticRules
}

// GetRule returns a rule by ID, or nil if not found.
func (r *Registry) GetRule(id string) Rule {
	return r.allRules[id]
}

// AllRules returns all registered rules.
func (r *Registry) AllRules() map[string]Rule {
	return r.allRules
}

// RunSyntaxRules runs all syntax rules against a parse result.
func (r *Registry) RunSyntaxRules(result *parser.ParseResult) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	for _, rule := range r.syntaxRules {
		diags = append(diags, rule.CheckSyntax(result)...)
	}
	return diags
}

// RunSemanticRules runs all semantic rules against parsed rules.
// Semantic rules can run even if KB is nil, as long as we have parsed rules from ANTLR.
func (r *Registry) RunSemanticRules(result *parser.ParseResult) []diagnostic.Diagnostic {
	// We need parsed rules to run semantic analysis
	if len(result.Rules) == 0 {
		return nil
	}

	var diags []diagnostic.Diagnostic
	for _, rule := range r.semanticRules {
		diags = append(diags, rule.CheckKnowledgeBase(result.File, result, result.KnowledgeBase)...)
	}
	return diags
}

// RegistryConfig holds configuration options for creating a registry.
type RegistryConfig struct {
	// NamingConvention specifies the naming convention for GRL007.
	// Valid values: "PascalCase", "camelCase", "snake_case", "kebab-case"
	NamingConvention string

	// MaxConditions specifies the max conditions for GRL006.
	MaxConditions int
}

// DefaultRegistry creates a registry with all built-in rules registered.
func DefaultRegistry() *Registry {
	return DefaultRegistryWithConfig(RegistryConfig{})
}

// DefaultRegistryWithConfig creates a registry with all built-in rules
// registered using the provided configuration.
func DefaultRegistryWithConfig(cfg RegistryConfig) *Registry {
	registry := NewRegistry()

	// Register syntax rules
	registry.RegisterSyntax(&SyntaxErrorRule{})

	// Register semantic rules (Phase 2)
	registry.RegisterSemantic(&MissingDescriptionRule{})
	registry.RegisterSemantic(&MissingSalienceRule{})
	registry.RegisterSemantic(&MissingRetractRule{})
	registry.RegisterSemantic(&DuplicateRuleRule{})

	// High complexity rule with configurable max conditions
	complexityRule := NewHighComplexityRule()
	if cfg.MaxConditions > 0 {
		complexityRule.MaxConditions = cfg.MaxConditions
	}
	registry.RegisterSemantic(complexityRule)

	// Register advanced semantic rules (Phase 4)
	// Naming convention rule with configurable convention
	if cfg.NamingConvention != "" {
		registry.RegisterSemantic(NewNamingConventionRuleWithConfig(cfg.NamingConvention))
	} else {
		registry.RegisterSemantic(NewNamingConventionRule())
	}

	registry.RegisterSemantic(&EmptyWhenRule{})
	registry.RegisterSemantic(&EmptyThenRule{})
	registry.RegisterSemantic(&UnusedVariableRule{})
	registry.RegisterSemantic(NewUndefinedVariableRule())
	registry.RegisterSemantic(&ConflictingRulesRule{})

	return registry
}
