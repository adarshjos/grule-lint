package lint_test

import (
	"fmt"

	"github.com/adarshjos/grule-lint/pkg/lint"
)

func ExampleNew() {
	// Create a linter with default configuration
	linter := lint.New()

	result := linter.LintString("example.grl", `
rule ProcessOrder "Process pending orders" salience 100 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
        Retract("ProcessOrder");
}
`)

	fmt.Printf("Found %d issues\n", result.Count())
	// Output: Found 0 issues
}

func ExampleNewWithConfig() {
	// Create configuration with snake_case naming convention
	cfg := lint.DefaultConfig()
	cfg.SetNamingConvention("snake_case")

	linter := lint.NewWithConfig(cfg)

	// PascalCase rule name will trigger GRL007 with snake_case convention
	result := linter.LintString("example.grl", `
rule MyRule "A rule" salience 10 {
    when true
    then Retract("MyRule");
}
`)

	for _, d := range result.All() {
		if d.RuleID() == "GRL007" {
			fmt.Println("Naming convention violation found")
		}
	}
	// Output: Naming convention violation found
}

func ExampleLinter_LintString() {
	linter := lint.New()

	// Lint GRL content from a string
	result := linter.LintString("rules.grl", `
rule ValidRule "A valid rule" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
        Retract("ValidRule");
}
`)

	if result.HasErrors() {
		fmt.Println("Errors found!")
	} else {
		fmt.Println("No errors found")
	}
	// Output: No errors found
}

func ExampleResult_Sorted() {
	linter := lint.New()

	result := linter.LintString("example.grl", `
rule RuleB salience 10 {
    when true
    then Log("b");
}
rule RuleA salience 10 {
    when true
    then Log("a");
}
`)

	// Get diagnostics sorted by file, line, column
	sorted := result.Sorted()
	fmt.Printf("Got %d diagnostics (sorted)\n", len(sorted))
	// Output: Got 7 diagnostics (sorted)
}

func ExampleAvailableRules() {
	// List all available lint rules
	rules := lint.AvailableRules()

	for _, r := range rules[:3] {
		fmt.Printf("%s: %s\n", r.ID, r.Name)
	}
	// Output:
	// GRL001: syntax-error
	// GRL002: missing-description
	// GRL003: missing-salience
}

func ExampleRuleByID() {
	// Get information about a specific rule
	rule := lint.RuleByID("GRL004")
	if rule != nil {
		fmt.Printf("%s: %s\n", rule.ID, rule.Description)
	}
	// Output: GRL004: Rule does not call Retract() which may cause infinite loops
}

func ExampleConfig_SetRule() {
	cfg := lint.DefaultConfig()

	// Configure rule severities (for use with reporters)
	cfg.SetRule("GRL003", "off")
	cfg.SetRule("GRL004", "error")

	// Check if a rule is enabled
	fmt.Printf("GRL003 enabled: %v\n", cfg.IsRuleEnabled("GRL003"))
	fmt.Printf("GRL004 enabled: %v\n", cfg.IsRuleEnabled("GRL004"))
	// Output:
	// GRL003 enabled: false
	// GRL004 enabled: true
}
