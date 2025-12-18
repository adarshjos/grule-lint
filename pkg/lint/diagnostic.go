package lint

import (
	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

// Severity represents the severity level of a diagnostic.
type Severity = diagnostic.Severity

// Severity levels from most to least severe.
const (
	SeverityError   = diagnostic.SeverityError
	SeverityWarning = diagnostic.SeverityWarning
	SeverityInfo    = diagnostic.SeverityInfo
	SeverityHint    = diagnostic.SeverityHint
)

// Position represents a location in a source file.
type Position = diagnostic.Position

// Range represents a span of text in a source file.
type Range = diagnostic.Range

// Fix represents a suggested fix for a diagnostic.
type Fix = diagnostic.Fix

// Edit represents a text edit.
type Edit = diagnostic.Edit

// Diagnostic represents a single lint issue found in a GRL file.
type Diagnostic struct {
	d diagnostic.Diagnostic
}

// File returns the file path where this diagnostic was found.
func (d Diagnostic) File() string {
	return d.d.File
}

// Range returns the source range of this diagnostic.
func (d Diagnostic) Range() Range {
	return d.d.Range
}

// Line returns the 1-based line number where this diagnostic starts.
func (d Diagnostic) Line() int {
	return d.d.Range.Start.Line
}

// Column returns the 1-based column number where this diagnostic starts.
func (d Diagnostic) Column() int {
	return d.d.Range.Start.Column
}

// RuleID returns the identifier of the lint rule (e.g., "GRL001").
func (d Diagnostic) RuleID() string {
	return d.d.RuleID
}

// RuleName returns the human-readable name of the lint rule.
func (d Diagnostic) RuleName() string {
	return d.d.RuleName
}

// Severity returns the severity level of this diagnostic.
func (d Diagnostic) Severity() Severity {
	return d.d.Severity
}

// Message returns the diagnostic message.
func (d Diagnostic) Message() string {
	return d.d.Message
}

// Fixes returns any suggested fixes for this diagnostic.
func (d Diagnostic) Fixes() []Fix {
	return d.d.Fixes
}

// Result is a collection of diagnostics from a lint operation.
type Result struct {
	ds *diagnostic.DiagnosticSet
}

// All returns all diagnostics in the result.
func (r *Result) All() []Diagnostic {
	internal := r.ds.All()
	result := make([]Diagnostic, len(internal))
	for i, d := range internal {
		result[i] = Diagnostic{d: d}
	}
	return result
}

// Sorted returns diagnostics sorted by file, line, then column.
func (r *Result) Sorted() []Diagnostic {
	internal := r.ds.Sorted()
	result := make([]Diagnostic, len(internal))
	for i, d := range internal {
		result[i] = Diagnostic{d: d}
	}
	return result
}

// Count returns the total number of diagnostics.
func (r *Result) Count() int {
	return r.ds.Count()
}

// CountBySeverity returns a map of severity levels to their counts.
func (r *Result) CountBySeverity() map[Severity]int {
	return r.ds.CountBySeverity()
}

// HasErrors returns true if the result contains any error-level diagnostics.
func (r *Result) HasErrors() bool {
	return r.ds.HasErrors()
}

// Errors returns the count of error-level diagnostics.
func (r *Result) Errors() int {
	counts := r.ds.CountBySeverity()
	return counts[SeverityError]
}

// Warnings returns the count of warning-level diagnostics.
func (r *Result) Warnings() int {
	counts := r.ds.CountBySeverity()
	return counts[SeverityWarning]
}

// FilterBySeverity returns diagnostics with severity at or above minSeverity.
func (r *Result) FilterBySeverity(minSeverity Severity) []Diagnostic {
	internal := r.ds.FilterBySeverity(minSeverity)
	result := make([]Diagnostic, len(internal))
	for i, d := range internal {
		result[i] = Diagnostic{d: d}
	}
	return result
}

// wrapDiagnosticSet wraps an internal DiagnosticSet as a public Result.
func wrapDiagnosticSet(ds *diagnostic.DiagnosticSet) *Result {
	if ds == nil {
		ds = diagnostic.NewDiagnosticSet()
	}
	return &Result{ds: ds}
}
