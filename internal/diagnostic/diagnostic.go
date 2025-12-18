// Package diagnostic provides types for representing lint diagnostics.
//
// A diagnostic represents a single issue found during linting, including
// its location, severity, and message. Diagnostics are collected into
// a DiagnosticSet for easy aggregation and filtering.
package diagnostic

import (
	"fmt"
	"sort"
)

// Severity represents the severity level of a diagnostic.
// Lower values indicate higher severity (Error < Warning < Info < Hint).
type Severity int

// Severity levels from most to least severe.
const (
	SeverityError   Severity = iota // SeverityError indicates a critical issue that should fail linting
	SeverityWarning                 // SeverityWarning indicates a potential problem
	SeverityInfo                    // SeverityInfo indicates informational messages
	SeverityHint                    // SeverityHint indicates suggestions for improvement
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	case SeverityHint:
		return "hint"
	default:
		return "unknown"
	}
}

// Position represents a location in a source file.
type Position struct {
	Line   int // Line is the 1-based line number
	Column int // Column is the 1-based column number
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Range represents a span of text in a source file.
type Range struct {
	Start Position // Start is the beginning of the range (inclusive)
	End   Position // End is the end of the range (exclusive)
}

// Fix represents a suggested fix for a diagnostic.
type Fix struct {
	Description string
	Edits       []Edit
}

// Edit represents a text edit.
type Edit struct {
	Range   Range
	NewText string
}

// Diagnostic represents a single lint issue.
type Diagnostic struct {
	File     string
	Range    Range
	RuleID   string
	RuleName string
	Severity Severity
	Message  string
	Fixes    []Fix // Optional suggested fixes
}

// DiagnosticSet is a collection of diagnostics with helper methods
// for aggregation, filtering, and sorting.
type DiagnosticSet struct {
	diagnostics []Diagnostic
}

// NewDiagnosticSet creates a new empty DiagnosticSet.
func NewDiagnosticSet() *DiagnosticSet {
	return &DiagnosticSet{
		diagnostics: make([]Diagnostic, 0),
	}
}

// Add appends a single diagnostic to the set.
func (ds *DiagnosticSet) Add(d Diagnostic) {
	ds.diagnostics = append(ds.diagnostics, d)
}

// AddAll appends multiple diagnostics to the set.
func (ds *DiagnosticSet) AddAll(diagnostics []Diagnostic) {
	ds.diagnostics = append(ds.diagnostics, diagnostics...)
}

// All returns all diagnostics in the set.
func (ds *DiagnosticSet) All() []Diagnostic {
	return ds.diagnostics
}

// Sorted returns a copy of diagnostics sorted by file, line, then column.
func (ds *DiagnosticSet) Sorted() []Diagnostic {
	sorted := make([]Diagnostic, len(ds.diagnostics))
	copy(sorted, ds.diagnostics)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].File != sorted[j].File {
			return sorted[i].File < sorted[j].File
		}
		if sorted[i].Range.Start.Line != sorted[j].Range.Start.Line {
			return sorted[i].Range.Start.Line < sorted[j].Range.Start.Line
		}
		return sorted[i].Range.Start.Column < sorted[j].Range.Start.Column
	})

	return sorted
}

// Count returns the total number of diagnostics in the set.
func (ds *DiagnosticSet) Count() int {
	return len(ds.diagnostics)
}

// CountBySeverity returns a map of severity levels to their counts.
func (ds *DiagnosticSet) CountBySeverity() map[Severity]int {
	counts := make(map[Severity]int)
	for _, d := range ds.diagnostics {
		counts[d.Severity]++
	}
	return counts
}

// HasErrors returns true if the set contains any error-level diagnostics.
func (ds *DiagnosticSet) HasErrors() bool {
	for _, d := range ds.diagnostics {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

// FilterBySeverity returns diagnostics with severity at or above minSeverity.
// Since lower Severity values are more severe, this returns diagnostics
// where d.Severity <= minSeverity.
func (ds *DiagnosticSet) FilterBySeverity(minSeverity Severity) []Diagnostic {
	var filtered []Diagnostic
	for _, d := range ds.diagnostics {
		if d.Severity <= minSeverity {
			filtered = append(filtered, d)
		}
	}
	return filtered
}
