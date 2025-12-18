package diagnostic

import "testing"

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev      Severity
		expected string
	}{
		{SeverityError, "error"},
		{SeverityWarning, "warning"},
		{SeverityInfo, "info"},
		{SeverityHint, "hint"},
		{Severity(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.expected {
			t.Errorf("Severity(%d).String() = %s, want %s", tt.sev, got, tt.expected)
		}
	}
}

func TestDiagnosticSet(t *testing.T) {
	ds := NewDiagnosticSet()

	if ds.Count() != 0 {
		t.Errorf("new set should be empty")
	}

	ds.Add(Diagnostic{RuleID: "GRL001", Severity: SeverityError})
	ds.Add(Diagnostic{RuleID: "GRL002", Severity: SeverityWarning})

	if ds.Count() != 2 {
		t.Errorf("expected 2, got %d", ds.Count())
	}

	if !ds.HasErrors() {
		t.Error("should have errors")
	}

	counts := ds.CountBySeverity()
	if counts[SeverityError] != 1 {
		t.Errorf("expected 1 error, got %d", counts[SeverityError])
	}
}

func TestDiagnosticSetSorted(t *testing.T) {
	ds := NewDiagnosticSet()
	ds.Add(Diagnostic{File: "b.grl", Range: Range{Start: Position{Line: 10}}})
	ds.Add(Diagnostic{File: "a.grl", Range: Range{Start: Position{Line: 5}}})
	ds.Add(Diagnostic{File: "a.grl", Range: Range{Start: Position{Line: 1}}})

	sorted := ds.Sorted()
	if sorted[0].File != "a.grl" || sorted[0].Range.Start.Line != 1 {
		t.Error("sorting failed")
	}
}

func TestFilterBySeverity(t *testing.T) {
	ds := NewDiagnosticSet()
	ds.Add(Diagnostic{Severity: SeverityError})
	ds.Add(Diagnostic{Severity: SeverityWarning})
	ds.Add(Diagnostic{Severity: SeverityInfo})

	filtered := ds.FilterBySeverity(SeverityWarning)
	if len(filtered) != 2 {
		t.Errorf("expected 2 (error+warning), got %d", len(filtered))
	}
}
