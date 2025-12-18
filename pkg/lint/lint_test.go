package lint_test

import (
	"testing"

	"github.com/adarshjos/grule-lint/pkg/lint"
)

func TestNew(t *testing.T) {
	linter := lint.New()
	if linter == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewWithConfig(t *testing.T) {
	cfg := lint.DefaultConfig()
	linter := lint.NewWithConfig(cfg)
	if linter == nil {
		t.Fatal("NewWithConfig() returned nil")
	}
}

func TestLintString_ValidGRL(t *testing.T) {
	linter := lint.New()
	result := linter.LintString("test.grl", `
rule TestRule "A test rule" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
        Retract("TestRule");
}
`)

	if result == nil {
		t.Fatal("LintString returned nil")
	}

	// Valid GRL should have no errors
	if result.HasErrors() {
		t.Errorf("Expected no errors, got %d", result.Errors())
	}
}

func TestLintString_SyntaxError(t *testing.T) {
	linter := lint.New()
	result := linter.LintString("test.grl", `
rule InvalidRule {
    when
    then
}
`)

	if result == nil {
		t.Fatal("LintString returned nil")
	}

	// Should have syntax errors
	if !result.HasErrors() {
		t.Error("Expected syntax errors but got none")
	}
}

func TestLintString_MissingRetract(t *testing.T) {
	linter := lint.New()
	result := linter.LintString("test.grl", `
rule NoRetract "Missing retract" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
}
`)

	if result == nil {
		t.Fatal("LintString returned nil")
	}

	// Should have a warning for missing Retract
	found := false
	for _, d := range result.All() {
		if d.RuleID() == "GRL004" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected GRL004 (missing-retract) warning")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := lint.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.MaxConditions() != 5 {
		t.Errorf("Expected MaxConditions=5, got %d", cfg.MaxConditions())
	}
}

func TestConfig_SetRule(t *testing.T) {
	cfg := lint.DefaultConfig()
	cfg.SetRule("GRL004", "off")

	if cfg.IsRuleEnabled("GRL004") {
		t.Error("Expected GRL004 to be disabled")
	}
}

func TestAvailableRules(t *testing.T) {
	rules := lint.AvailableRules()
	if len(rules) == 0 {
		t.Fatal("No rules available")
	}

	// Check for expected rules
	expectedIDs := []string{"GRL001", "GRL002", "GRL003", "GRL004", "GRL005"}
	for _, id := range expectedIDs {
		found := false
		for _, r := range rules {
			if r.ID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule %s not found", id)
		}
	}
}

func TestRuleByID(t *testing.T) {
	rule := lint.RuleByID("GRL001")
	if rule == nil {
		t.Fatal("RuleByID(GRL001) returned nil")
	}

	if rule.Name != "syntax-error" {
		t.Errorf("Expected name 'syntax-error', got '%s'", rule.Name)
	}

	// Test non-existent rule
	if lint.RuleByID("INVALID") != nil {
		t.Error("Expected nil for invalid rule ID")
	}
}

func TestResult_Sorted(t *testing.T) {
	linter := lint.New()
	result := linter.LintString("test.grl", `
rule Second "desc" salience 10 {
    when true
    then Log("x");
}
rule First "desc" salience 10 {
    when true
    then Log("y");
}
`)

	sorted := result.Sorted()
	if len(sorted) == 0 {
		return // No diagnostics to sort
	}

	// Verify sorted by line
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Line() < sorted[i-1].Line() {
			t.Error("Diagnostics not sorted by line")
		}
	}
}

func TestDiagnostic_Accessors(t *testing.T) {
	linter := lint.New()
	result := linter.LintString("test.grl", `
rule Bad {
    when
    then
}
`)

	if result.Count() == 0 {
		t.Skip("No diagnostics to test")
	}

	d := result.All()[0]

	// Test all accessors don't panic
	_ = d.File()
	_ = d.Line()
	_ = d.Column()
	_ = d.Range()
	_ = d.RuleID()
	_ = d.RuleName()
	_ = d.Severity()
	_ = d.Message()
	_ = d.Fixes()
}

func TestNamingConvention_PascalCase(t *testing.T) {
	cfg := lint.DefaultConfig()
	cfg.SetNamingConvention("PascalCase")
	linter := lint.NewWithConfig(cfg)

	// snake_case name should trigger GRL007
	result := linter.LintString("test.grl", `
rule my_rule "desc" salience 10 {
    when true
    then Retract("my_rule");
}
`)

	found := false
	for _, d := range result.All() {
		if d.RuleID() == "GRL007" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected GRL007 for snake_case rule name with PascalCase convention")
	}
}

func TestNamingConvention_SnakeCase(t *testing.T) {
	cfg := lint.DefaultConfig()
	cfg.SetNamingConvention("snake_case")
	linter := lint.NewWithConfig(cfg)

	// PascalCase name should trigger GRL007
	result := linter.LintString("test.grl", `
rule MyRule "desc" salience 10 {
    when true
    then Retract("MyRule");
}
`)

	found := false
	for _, d := range result.All() {
		if d.RuleID() == "GRL007" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected GRL007 for PascalCase rule name with snake_case convention")
	}

	// snake_case name should NOT trigger GRL007
	result2 := linter.LintString("test.grl", `
rule my_rule "desc" salience 10 {
    when true
    then Retract("my_rule");
}
`)

	for _, d := range result2.All() {
		if d.RuleID() == "GRL007" {
			t.Error("Did not expect GRL007 for snake_case rule name with snake_case convention")
		}
	}
}
