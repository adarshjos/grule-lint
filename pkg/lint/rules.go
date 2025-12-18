package lint

// RuleInfo describes a lint rule.
type RuleInfo struct {
	ID          string   // Unique identifier (e.g., "GRL001")
	Name        string   // Human-readable name (e.g., "syntax-error")
	Description string   // Brief description of what the rule checks
	Severity    Severity // Default severity level
}

// AvailableRules returns information about all available lint rules.
func AvailableRules() []RuleInfo {
	return []RuleInfo{
		{
			ID:          "GRL001",
			Name:        "syntax-error",
			Description: "Reports syntax errors in GRL files",
			Severity:    SeverityError,
		},
		{
			ID:          "GRL002",
			Name:        "missing-description",
			Description: "Rule is missing a description string",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL003",
			Name:        "missing-salience",
			Description: "Rule is missing an explicit salience value",
			Severity:    SeverityInfo,
		},
		{
			ID:          "GRL004",
			Name:        "missing-retract",
			Description: "Rule does not call Retract() which may cause infinite loops",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL005",
			Name:        "duplicate-rule",
			Description: "Multiple rules have the same name",
			Severity:    SeverityError,
		},
		{
			ID:          "GRL006",
			Name:        "high-complexity",
			Description: "Rule has too many conditions in the when clause",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL007",
			Name:        "naming-convention",
			Description: "Rule name does not follow naming conventions",
			Severity:    SeverityInfo,
		},
		{
			ID:          "GRL008",
			Name:        "unused-variable",
			Description: "Variable is assigned but never used",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL009",
			Name:        "empty-when",
			Description: "Rule has an empty when clause",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL010",
			Name:        "empty-then",
			Description: "Rule has an empty then clause",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL011",
			Name:        "constant-condition",
			Description: "When clause contains a constant condition (always true/false)",
			Severity:    SeverityWarning,
		},
		{
			ID:          "GRL012",
			Name:        "unreachable-rule",
			Description: "Rule is unreachable due to conflicting conditions",
			Severity:    SeverityWarning,
		},
	}
}

// RuleByID returns information about a specific rule, or nil if not found.
func RuleByID(id string) *RuleInfo {
	for _, r := range AvailableRules() {
		if r.ID == id {
			return &r
		}
	}
	return nil
}
