package test

import (
	"testing"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/linter"
)

// TestLinter_SyntaxErrors tests that syntax errors are detected
func TestLinter_SyntaxErrors(t *testing.T) {
	l := linter.New()

	content := `
rule BadRule "Missing when" {
    then
        Log("Bad");
}
`
	ds := l.LintString("test.grl", content)

	if ds.Count() == 0 {
		t.Fatal("Expected syntax errors")
	}

	// Should have GRL001 error
	found := false
	for _, d := range ds.All() {
		if d.RuleID == "GRL001" && d.Severity == diagnostic.SeverityError {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected GRL001 syntax error")
	}
}

// TestLinter_MissingDescription tests GRL002
func TestLinter_MissingDescription(t *testing.T) {
	l := linter.New()

	content := `
rule NoDescription salience 1 {
    when Order.Status == "pending"
    then Retract("NoDescription");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL002")
	if !found {
		t.Error("Expected GRL002 missing-description warning")
	}
}

// TestLinter_MissingSalience tests GRL003
func TestLinter_MissingSalience(t *testing.T) {
	l := linter.New()

	content := `
rule NoSalience "Test" {
    when Order.Status == "pending"
    then Retract("NoSalience");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL003")
	if !found {
		t.Error("Expected GRL003 missing-salience info")
	}
}

// TestLinter_MissingRetract tests GRL004
func TestLinter_MissingRetract(t *testing.T) {
	l := linter.New()

	content := `
rule NoRetract "Test" salience 1 {
    when Order.Status == "pending"
    then Log("No retract");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL004")
	if !found {
		t.Error("Expected GRL004 missing-retract warning")
	}
}

// TestLinter_DuplicateRule tests GRL005
func TestLinter_DuplicateRule(t *testing.T) {
	l := linter.New()

	content := `
rule DuplicateRule "First" salience 1 {
    when Order.Status == "first"
    then Retract("DuplicateRule");
}

rule DuplicateRule "Second" salience 2 {
    when Order.Status == "second"
    then Retract("DuplicateRule");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL005")
	if !found {
		t.Error("Expected GRL005 duplicate-rule error")
	}
}

// TestLinter_HighComplexity tests GRL006
func TestLinter_HighComplexity(t *testing.T) {
	l := linter.New()

	content := `
rule ComplexRule "Complex" salience 1 {
    when
        Order.Status == "pending" &&
        Order.Total > 100 &&
        Order.Priority == "high" &&
        Customer.Type == "premium" &&
        Customer.Age > 18 &&
        Product.InStock == true
    then
        Retract("ComplexRule");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL006")
	if !found {
		t.Error("Expected GRL006 high-complexity warning")
	}
}

// TestLinter_NamingConvention tests GRL007
func TestLinter_NamingConvention(t *testing.T) {
	l := linter.New()

	// snake_case should fail PascalCase convention
	content := `
rule snake_case_rule "Test" salience 1 {
    when Order.Status == "pending"
    then Retract("snake_case_rule");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL007")
	if !found {
		t.Error("Expected GRL007 naming-convention warning")
	}
}

// TestLinter_EmptyWhen tests GRL010
// Detects 'when true' as an always-true condition
func TestLinter_EmptyWhen(t *testing.T) {
	l := linter.New()

	// Using 'true' literal - this is an always-true condition with no actual comparison
	content := `
rule AlwaysTrueRule "Always true when clause" salience 1 {
    when
        true
    then
        Retract("AlwaysTrueRule");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL010")
	if !found {
		t.Error("Expected GRL010 empty-when warning for 'when true' clause")
	}
}

// TestLinter_EmptyThen tests GRL011
// Note: GRL syntax requires at least one statement in then clause, so we test with assignment only
func TestLinter_EmptyThen(t *testing.T) {
	l := linter.New()

	// A rule with only an assignment (no function calls like Retract/Log)
	content := `
rule AssignOnlyRule "Assignment only in then" salience 1 {
    when Order.Status == "pending"
    then
        Order.Processed = true;
}
`
	ds := l.LintString("test.grl", content)

	// GRL011 checks for FunctionCalls == 0 && ThenActionCount == 0
	// Since we have an assignment, ThenActionCount > 0, so this won't trigger
	// This validates GRL011 doesn't false-positive on assignment-only rules
	found := hasRuleID(ds, "GRL011")
	if found {
		t.Error("GRL011 should not trigger for rules with assignments")
	}
}

// TestLinter_ConflictingRules tests GRL012
func TestLinter_ConflictingRules(t *testing.T) {
	l := linter.New()

	content := `
rule ConflictA "First" {
    when Order.Total > 100
    then Retract("ConflictA");
}

rule ConflictB "Second" {
    when Order.Total > 100
    then Retract("ConflictB");
}
`
	ds := l.LintString("test.grl", content)

	found := hasRuleID(ds, "GRL012")
	if !found {
		t.Error("Expected GRL012 conflicting-rules warning")
	}
}

// TestLinter_ValidRule tests that a valid rule produces no errors
func TestLinter_ValidRule(t *testing.T) {
	l := linter.New()

	content := `
rule ValidRule "A properly formatted rule" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processed";
        Retract("ValidRule");
}
`
	ds := l.LintString("test.grl", content)

	// Should only have GRL003 (missing-salience is info, but salience is present)
	// and possibly GRL007 if naming doesn't match
	errorCount := 0
	for _, d := range ds.All() {
		if d.Severity == diagnostic.SeverityError {
			errorCount++
			t.Errorf("Unexpected error: %s - %s", d.RuleID, d.Message)
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors for valid rule, got %d", errorCount)
	}
}

// TestLinter_MultipleIssues tests multiple rules firing
func TestLinter_MultipleIssues(t *testing.T) {
	l := linter.New()

	// This rule has multiple issues
	content := `
rule bad_rule salience 1 {
    when Order.Status == "pending"
    then Log("Missing retract and bad name");
}
`
	ds := l.LintString("test.grl", content)

	// Should have:
	// - GRL002 (missing description)
	// - GRL004 (missing retract)
	// - GRL007 (naming convention)

	expectedRules := []string{"GRL002", "GRL004", "GRL007"}
	for _, ruleID := range expectedRules {
		if !hasRuleID(ds, ruleID) {
			t.Errorf("Expected %s diagnostic", ruleID)
		}
	}
}

// TestLinter_FileWithMultipleRules tests linting multiple rules in one file
func TestLinter_FileWithMultipleRules(t *testing.T) {
	l := linter.New()

	content := `
rule FirstRule "First rule" salience 1 {
    when Order.Status == "first"
    then Retract("FirstRule");
}

rule SecondRule "Second rule" salience 2 {
    when Order.Status == "second"
    then Retract("SecondRule");
}

rule ThirdRule "Third rule" salience 3 {
    when Order.Status == "third"
    then Retract("ThirdRule");
}
`
	ds := l.LintString("test.grl", content)

	// All rules are valid - should have no errors
	if ds.HasErrors() {
		for _, d := range ds.All() {
			if d.Severity == diagnostic.SeverityError {
				t.Errorf("Unexpected error: %s - %s", d.RuleID, d.Message)
			}
		}
	}
}

// TestLinter_RealWorldExample tests a realistic rule file
func TestLinter_RealWorldExample(t *testing.T) {
	l := linter.New()

	content := `
// Order processing rules

rule ProcessPendingOrder "Process orders that are pending" salience 10 {
    when
        Order.Status == "pending" &&
        Order.PaymentVerified == true
    then
        Order.Status = "processing";
        Log("Order moved to processing");
        Retract("ProcessPendingOrder");
}

rule ApplyPremiumDiscount "Apply discount for premium customers" salience 5 {
    when
        Customer.Type == "premium" &&
        Order.Total > 100
    then
        Order.Discount = Order.Total * 0.1;
        Retract("ApplyPremiumDiscount");
}

rule FlagHighValueOrder "Flag orders over 1000" salience 1 {
    when
        Order.Total > 1000
    then
        Order.RequiresReview = true;
        Retract("FlagHighValueOrder");
}
`
	ds := l.LintString("test.grl", content)

	// Should have no errors
	if ds.HasErrors() {
		for _, d := range ds.All() {
			if d.Severity == diagnostic.SeverityError {
				t.Errorf("Unexpected error: %s at %d:%d - %s",
					d.RuleID, d.Range.Start.Line, d.Range.Start.Column, d.Message)
			}
		}
		t.Fatal("Expected no errors for well-formed rules")
	}
}

// TestLinter_PositionAccuracy tests that positions are reported correctly
func TestLinter_PositionAccuracy(t *testing.T) {
	l := linter.New()

	content := `rule BadName salience 1 {
    when Order.Status == "pending"
    then Retract("BadName");
}`
	ds := l.LintString("test.grl", content)

	// GRL002 should be at line 1 (rule start)
	for _, d := range ds.All() {
		if d.RuleID == "GRL002" {
			if d.Range.Start.Line != 1 {
				t.Errorf("Expected GRL002 at line 1, got line %d", d.Range.Start.Line)
			}
		}
	}
}
