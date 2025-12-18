package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

// TextReporter formats diagnostics as human-readable text output.
type TextReporter struct {
	writer    io.Writer
	colorized bool
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// NewTextReporter creates a new TextReporter.
func NewTextReporter(w io.Writer, colorized bool) *TextReporter {
	return &TextReporter{
		writer:    w,
		colorized: colorized,
	}
}

// Report outputs diagnostics in text format.
func (r *TextReporter) Report(diagnostics *diagnostic.DiagnosticSet) error {
	sorted := diagnostics.Sorted()

	for _, d := range sorted {
		line := r.formatDiagnostic(d)
		if _, err := fmt.Fprintln(r.writer, line); err != nil {
			return fmt.Errorf("writing diagnostic: %w", err)
		}
	}

	// Print summary
	if diagnostics.Count() > 0 {
		if _, err := fmt.Fprintln(r.writer); err != nil {
			return fmt.Errorf("writing newline: %w", err)
		}
		if _, err := fmt.Fprintln(r.writer, r.formatSummary(diagnostics)); err != nil {
			return fmt.Errorf("writing summary: %w", err)
		}
	}

	return nil
}

// formatDiagnostic formats a single diagnostic.
func (r *TextReporter) formatDiagnostic(d diagnostic.Diagnostic) string {
	// Format: file:line:col: CODE [severity] message
	location := fmt.Sprintf("%s:%d:%d:",
		d.File,
		d.Range.Start.Line,
		d.Range.Start.Column,
	)

	severity := fmt.Sprintf("[%s]", d.Severity.String())

	if r.colorized {
		location = colorGray + location + colorReset
		severity = r.colorSeverity(d.Severity, severity)
		return fmt.Sprintf("%s %s%s%s %s %s",
			location,
			colorBold,
			d.RuleID,
			colorReset,
			severity,
			d.Message,
		)
	}

	return fmt.Sprintf("%s %s %s %s", location, d.RuleID, severity, d.Message)
}

// colorSeverity applies color to severity text.
func (r *TextReporter) colorSeverity(s diagnostic.Severity, text string) string {
	switch s {
	case diagnostic.SeverityError:
		return colorRed + text + colorReset
	case diagnostic.SeverityWarning:
		return colorYellow + text + colorReset
	case diagnostic.SeverityInfo:
		return colorCyan + text + colorReset
	default:
		return text
	}
}

// formatSummary formats the summary line.
func (r *TextReporter) formatSummary(diagnostics *diagnostic.DiagnosticSet) string {
	counts := diagnostics.CountBySeverity()

	var parts []string
	if errors := counts[diagnostic.SeverityError]; errors > 0 {
		part := fmt.Sprintf("%d error", errors)
		if errors != 1 {
			part += "s"
		}
		if r.colorized {
			part = colorRed + part + colorReset
		}
		parts = append(parts, part)
	}

	if warnings := counts[diagnostic.SeverityWarning]; warnings > 0 {
		part := fmt.Sprintf("%d warning", warnings)
		if warnings != 1 {
			part += "s"
		}
		if r.colorized {
			part = colorYellow + part + colorReset
		}
		parts = append(parts, part)
	}

	if infos := counts[diagnostic.SeverityInfo]; infos > 0 {
		part := fmt.Sprintf("%d info", infos)
		parts = append(parts, part)
	}

	total := diagnostics.Count()
	issue := "issue"
	if total != 1 {
		issue = "issues"
	}

	if len(parts) == 0 {
		return fmt.Sprintf("Found %d %s", total, issue)
	}

	return fmt.Sprintf("Found %d %s (%s)", total, issue, strings.Join(parts, ", "))
}
