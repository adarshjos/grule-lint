package test

import (
	"testing"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

// hasRuleID checks if a diagnostic with the given rule ID exists in the set.
func hasRuleID(ds *diagnostic.DiagnosticSet, ruleID string) bool {
	for _, d := range ds.All() {
		if d.RuleID == ruleID {
			return true
		}
	}
	return false
}

// countRuleID counts how many diagnostics have the given rule ID.
func countRuleID(ds *diagnostic.DiagnosticSet, ruleID string) int {
	count := 0
	for _, d := range ds.All() {
		if d.RuleID == ruleID {
			count++
		}
	}
	return count
}

// hasSeverity checks if a diagnostic with the given severity exists.
func hasSeverity(ds *diagnostic.DiagnosticSet, severity diagnostic.Severity) bool {
	for _, d := range ds.All() {
		if d.Severity == severity {
			return true
		}
	}
	return false
}

// countSeverity counts diagnostics with the given severity.
func countSeverity(ds *diagnostic.DiagnosticSet, severity diagnostic.Severity) int {
	count := 0
	for _, d := range ds.All() {
		if d.Severity == severity {
			count++
		}
	}
	return count
}

// assertHasRuleID fails the test if the rule ID is not present.
func assertHasRuleID(t *testing.T, ds *diagnostic.DiagnosticSet, ruleID string) {
	t.Helper()
	if !hasRuleID(ds, ruleID) {
		t.Errorf("Expected diagnostic with rule ID %s, but not found", ruleID)
	}
}

// assertNoRuleID fails the test if the rule ID is present.
func assertNoRuleID(t *testing.T, ds *diagnostic.DiagnosticSet, ruleID string) {
	t.Helper()
	if hasRuleID(ds, ruleID) {
		t.Errorf("Did not expect diagnostic with rule ID %s, but found", ruleID)
	}
}

// assertNoErrors fails the test if there are any error-level diagnostics.
func assertNoErrors(t *testing.T, ds *diagnostic.DiagnosticSet) {
	t.Helper()
	if ds.HasErrors() {
		for _, d := range ds.All() {
			if d.Severity == diagnostic.SeverityError {
				t.Errorf("Unexpected error: %s at %d:%d - %s",
					d.RuleID, d.Range.Start.Line, d.Range.Start.Column, d.Message)
			}
		}
	}
}

// getDiagnosticByRuleID returns the first diagnostic with the given rule ID.
func getDiagnosticByRuleID(ds *diagnostic.DiagnosticSet, ruleID string) *diagnostic.Diagnostic {
	for _, d := range ds.All() {
		if d.RuleID == ruleID {
			return &d
		}
	}
	return nil
}
