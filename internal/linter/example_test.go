package linter_test

import (
	"fmt"

	"github.com/adarshjos/grule-lint/internal/linter"
)

func ExampleLinter_LintString() {
	l := linter.New()

	result := l.LintString("example.grl", `
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

func ExampleLinter_LintString_withErrors() {
	l := linter.New()

	// Rule missing Retract() call
	result := l.LintString("example.grl", `
rule NoRetract "Missing retract" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
}
`)

	for _, d := range result.All() {
		if d.RuleID == "GRL004" {
			fmt.Printf("Rule: %s - %s\n", d.RuleID, d.RuleName)
		}
	}
	// Output: Rule: GRL004 - missing-retract
}
