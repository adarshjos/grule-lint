package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

func TestTextReporter(t *testing.T) {
	ds := diagnostic.NewDiagnosticSet()
	ds.Add(diagnostic.Diagnostic{
		File:     "test.grl",
		Range:    diagnostic.Range{Start: diagnostic.Position{Line: 10, Column: 5}},
		RuleID:   "GRL001",
		Severity: diagnostic.SeverityError,
		Message:  "syntax error",
	})

	var buf bytes.Buffer
	reporter := NewTextReporter(&buf, false)
	reporter.Report(ds)

	output := buf.String()
	if !strings.Contains(output, "test.grl:10:5") {
		t.Errorf("expected file:line:col, got: %s", output)
	}
	if !strings.Contains(output, "GRL001") {
		t.Error("expected rule ID in output")
	}
	if !strings.Contains(output, "1 error") {
		t.Error("expected summary line")
	}
}

func TestTextReporterEmpty(t *testing.T) {
	ds := diagnostic.NewDiagnosticSet()

	var buf bytes.Buffer
	reporter := NewTextReporter(&buf, false)
	reporter.Report(ds)

	// Empty diagnostic set produces no output (summary only shown when count > 0)
	if buf.String() != "" {
		t.Errorf("expected empty output for no issues, got: %s", buf.String())
	}
}
