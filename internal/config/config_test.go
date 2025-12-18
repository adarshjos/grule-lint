package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.Complexity.MaxConditions != 5 {
		t.Errorf("expected MaxConditions=5, got %d", cfg.Complexity.MaxConditions)
	}
}

func TestLoad(t *testing.T) {
	// Create temp config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".grl-lint.yaml")
	content := []byte(`
rules:
  GRL001: error
  GRL002: off
complexity:
  max_conditions: 10
`)
	if err := os.WriteFile(configPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Rules["GRL001"] != "error" {
		t.Errorf("expected GRL001=error, got %s", cfg.Rules["GRL001"])
	}
	if cfg.Complexity.MaxConditions != 10 {
		t.Errorf("expected MaxConditions=10, got %d", cfg.Complexity.MaxConditions)
	}
}

func TestGetRuleSeverity(t *testing.T) {
	cfg := &Config{Rules: map[string]string{
		"GRL001": "error",
		"GRL002": "warning",
		"GRL003": "off",
	}}

	tests := []struct {
		ruleID   string
		expected *diagnostic.Severity
	}{
		{"GRL001", ptr(diagnostic.SeverityError)},
		{"GRL002", ptr(diagnostic.SeverityWarning)},
		{"GRL003", nil},
		{"GRL999", ptr(diagnostic.SeverityWarning)}, // default
	}

	for _, tt := range tests {
		sev := cfg.GetRuleSeverity(tt.ruleID, diagnostic.SeverityWarning)
		if tt.expected == nil && sev != nil {
			t.Errorf("%s: expected nil, got %v", tt.ruleID, *sev)
		} else if tt.expected != nil && (sev == nil || *sev != *tt.expected) {
			t.Errorf("%s: expected %v, got %v", tt.ruleID, *tt.expected, sev)
		}
	}
}

func TestIsRuleEnabled(t *testing.T) {
	cfg := &Config{Rules: map[string]string{
		"GRL001": "error",
		"GRL002": "off",
		"GRL003": "disabled",
	}}

	if !cfg.IsRuleEnabled("GRL001") {
		t.Error("GRL001 should be enabled")
	}
	if cfg.IsRuleEnabled("GRL002") {
		t.Error("GRL002 should be disabled")
	}
	if cfg.IsRuleEnabled("GRL003") {
		t.Error("GRL003 should be disabled")
	}
	if !cfg.IsRuleEnabled("GRL999") {
		t.Error("unknown rule should be enabled by default")
	}
}

func TestShouldExclude(t *testing.T) {
	cfg := &Config{Exclude: []string{"**/vendor/**", "*.tmp"}}

	tests := []struct {
		path     string
		excluded bool
	}{
		{"vendor/lib.grl", true}, // ** pattern matches
		{"file.tmp", true},
		{"rules/order.grl", false},
	}

	for _, tt := range tests {
		if got := cfg.ShouldExclude(tt.path); got != tt.excluded {
			t.Errorf("ShouldExclude(%s): expected %v, got %v", tt.path, tt.excluded, got)
		}
	}
}

func ptr(s diagnostic.Severity) *diagnostic.Severity {
	return &s
}
