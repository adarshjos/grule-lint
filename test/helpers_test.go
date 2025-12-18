package test

import (
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
