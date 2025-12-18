package test

import (
	"testing"

	"github.com/adarshjos/grule-lint/internal/linter"
)

// TestRules_Detection uses table-driven tests to verify each rule detects issues correctly.
func TestRules_Detection(t *testing.T) {
	tests := []struct {
		name        string
		grl         string
		expectRule  string // Expected rule ID to be triggered
		shouldExist bool   // Whether the rule should be triggered
	}{
		{
			name: "GRL001_SyntaxError",
			grl: `
rule BadSyntax "Missing when" {
    then Log("bad");
}`,
			expectRule:  "GRL001",
			shouldExist: true,
		},
		{
			name: "GRL002_MissingDescription",
			grl: `
rule NoDesc salience 1 {
    when Order.Status == "pending"
    then Retract("NoDesc");
}`,
			expectRule:  "GRL002",
			shouldExist: true,
		},
		{
			name: "GRL002_WithDescription_NoTrigger",
			grl: `
rule HasDesc "This is a description" salience 1 {
    when Order.Status == "pending"
    then Retract("HasDesc");
}`,
			expectRule:  "GRL002",
			shouldExist: false,
		},
		{
			name: "GRL003_MissingSalience",
			grl: `
rule NoSalience "Test" {
    when Order.Status == "pending"
    then Retract("NoSalience");
}`,
			expectRule:  "GRL003",
			shouldExist: true,
		},
		{
			name: "GRL003_WithSalience_NoTrigger",
			grl: `
rule HasSalience "Test" salience 10 {
    when Order.Status == "pending"
    then Retract("HasSalience");
}`,
			expectRule:  "GRL003",
			shouldExist: false,
		},
		{
			name: "GRL004_MissingRetract",
			grl: `
rule NoRetract "Test" salience 1 {
    when Order.Status == "pending"
    then Log("no retract");
}`,
			expectRule:  "GRL004",
			shouldExist: true,
		},
		{
			name: "GRL004_WithRetract_NoTrigger",
			grl: `
rule HasRetract "Test" salience 1 {
    when Order.Status == "pending"
    then Retract("HasRetract");
}`,
			expectRule:  "GRL004",
			shouldExist: false,
		},
		{
			name: "GRL006_HighComplexity",
			grl: `
rule Complex "Complex" salience 1 {
    when
        A == 1 && B == 2 && C == 3 && D == 4 && E == 5 && F == 6
    then Retract("Complex");
}`,
			expectRule:  "GRL006",
			shouldExist: true,
		},
		{
			name: "GRL006_LowComplexity_NoTrigger",
			grl: `
rule Simple "Simple" salience 1 {
    when A == 1 && B == 2
    then Retract("Simple");
}`,
			expectRule:  "GRL006",
			shouldExist: false,
		},
		{
			name: "GRL007_NamingConvention_SnakeCase",
			grl: `
rule snake_case "Test" salience 1 {
    when Order.Status == "pending"
    then Retract("snake_case");
}`,
			expectRule:  "GRL007",
			shouldExist: true,
		},
		{
			name: "GRL007_NamingConvention_PascalCase_NoTrigger",
			grl: `
rule PascalCase "Test" salience 1 {
    when Order.Status == "pending"
    then Retract("PascalCase");
}`,
			expectRule:  "GRL007",
			shouldExist: false,
		},
		{
			name: "GRL010_EmptyWhen_TrueLiteral",
			grl: `
rule AlwaysTrue "Test" salience 1 {
    when true
    then Retract("AlwaysTrue");
}`,
			expectRule:  "GRL010",
			shouldExist: true,
		},
	}

	l := linter.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := l.LintString("test.grl", tt.grl)
			exists := hasRuleID(ds, tt.expectRule)

			if tt.shouldExist && !exists {
				t.Errorf("Expected %s to be triggered, but it wasn't. Diagnostics: %v",
					tt.expectRule, ds.All())
			}
			if !tt.shouldExist && exists {
				t.Errorf("Did not expect %s to be triggered, but it was",
					tt.expectRule)
			}
		})
	}
}

// TestRules_MultipleIssues tests that multiple rules fire on the same input.
func TestRules_MultipleIssues(t *testing.T) {
	tests := []struct {
		name         string
		grl          string
		expectRules  []string
		excludeRules []string
	}{
		{
			name: "BadRule_MultipleViolations",
			grl: `
rule bad_rule salience 1 {
    when Order.Status == "pending"
    then Log("test");
}`,
			expectRules:  []string{"GRL002", "GRL004", "GRL007"},
			excludeRules: []string{"GRL001"},
		},
		{
			name: "ValidRule_NoViolations",
			grl: `
rule ValidRule "Description" salience 10 {
    when Order.Status == "pending"
    then Retract("ValidRule");
}`,
			expectRules:  []string{},
			excludeRules: []string{"GRL001", "GRL002", "GRL003", "GRL004", "GRL007"},
		},
	}

	l := linter.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := l.LintString("test.grl", tt.grl)

			for _, ruleID := range tt.expectRules {
				if !hasRuleID(ds, ruleID) {
					t.Errorf("Expected %s to be triggered", ruleID)
				}
			}

			for _, ruleID := range tt.excludeRules {
				if hasRuleID(ds, ruleID) {
					t.Errorf("Did not expect %s to be triggered", ruleID)
				}
			}
		})
	}
}
