package parser

import (
	"fmt"
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/hyperjumptech/grule-rule-engine/antlr/parser/grulev3"
	gruleAst "github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

// DefaultMaxFileSize is the default maximum file size for parsing (10MB).
const DefaultMaxFileSize = 10 * 1024 * 1024

// ParseError represents a syntax error from the Grule parser.
type ParseError struct {
	Line    int
	Column  int
	Message string
}

// ParseResult contains the result of parsing a GRL file.
type ParseResult struct {
	// KnowledgeBase is the parsed AST (nil if parsing failed)
	KnowledgeBase *gruleAst.KnowledgeBase

	// Rules contains parsed rule information with accurate positions
	Rules []RuleInfo

	// Errors contains any parse errors encountered
	Errors []ParseError

	// Source is the original source content
	Source string

	// File is the file path
	File string
}

// Success returns true if parsing succeeded without errors.
func (r *ParseResult) Success() bool {
	return len(r.Errors) == 0 && r.KnowledgeBase != nil
}

// Parser wraps Grule's ANTLR parser for parsing GRL files.
type Parser struct {
	// MaxFileSize is the maximum file size allowed for parsing.
	// Set to 0 to use DefaultMaxFileSize. Set to -1 for unlimited.
	MaxFileSize int64
}

// NewParser creates a new Parser instance with default settings.
func NewParser() *Parser {
	return &Parser{
		MaxFileSize: DefaultMaxFileSize,
	}
}

// ParseString parses GRL content from a string.
func (p *Parser) ParseString(file, content string) *ParseResult {
	result := &ParseResult{
		Source: content,
		File:   file,
	}

	rules, errors := p.parseWithANTLR(content)
	result.Rules = rules
	result.Errors = errors

	if len(errors) > 0 {
		return result
	}

	// Build KnowledgeBase for semantic analysis (may fail for duplicate rules)
	lib := gruleAst.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	resource := pkg.NewBytesResource([]byte(content))

	if err := rb.BuildRuleFromResource("lint", "1.0.0", resource); err != nil {
		return result
	}

	result.KnowledgeBase = lib.GetKnowledgeBase("lint", "1.0.0")
	return result
}

// ParseFile parses a GRL file from disk.
func (p *Parser) ParseFile(file string) (*ParseResult, error) {
	// Check file size before reading
	maxSize := p.MaxFileSize
	if maxSize == 0 {
		maxSize = DefaultMaxFileSize
	}

	if maxSize > 0 {
		info, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("accessing GRL file %s: %w", file, err)
		}
		if info.Size() > maxSize {
			return nil, fmt.Errorf("GRL file %s exceeds maximum size limit (%d bytes > %d bytes)", file, info.Size(), maxSize)
		}
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading GRL file %s: %w", file, err)
	}

	result := p.ParseString(file, string(content))
	return result, nil
}

func (p *Parser) parseWithANTLR(content string) ([]RuleInfo, []ParseError) {
	input := antlr.NewInputStream(content)
	lexer := grulev3.Newgrulev3Lexer(input)

	errorReporter := &pkg.GruleErrorReporter{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorReporter)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := grulev3.Newgrulev3Parser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorReporter)
	parser.BuildParseTrees = true

	listener := NewLintListener()
	tree := parser.Grl()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.Rules, convertGruleErrors(errorReporter.Errors)
}

// convertGruleErrors parses "grl error on LINE:COL MSG" format.
func convertGruleErrors(gruleErrors []error) []ParseError {
	var errors []ParseError

	for _, err := range gruleErrors {
		pe := ParseError{Line: 1, Column: 1, Message: err.Error()}

		var line, col int
		if n, _ := fmt.Sscanf(err.Error(), "grl error on %d:%d", &line, &col); n == 2 {
			pe.Line = line
			pe.Column = col + 1
			prefix := fmt.Sprintf("grl error on %d:%d ", line, col)
			if len(err.Error()) > len(prefix) {
				pe.Message = err.Error()[len(prefix):]
			}
		}
		errors = append(errors, pe)
	}
	return errors
}

// ToDiagnostics converts parse errors to diagnostics.
func (r *ParseResult) ToDiagnostics() []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic

	for _, e := range r.Errors {
		diags = append(diags, diagnostic.Diagnostic{
			File: r.File,
			Range: diagnostic.Range{
				Start: diagnostic.Position{Line: e.Line, Column: e.Column},
				End:   diagnostic.Position{Line: e.Line, Column: e.Column},
			},
			RuleID:   "GRL001",
			RuleName: "syntax-error",
			Severity: diagnostic.SeverityError,
			Message:  e.Message,
		})
	}

	return diags
}

// GetRulePosition returns the position of a rule by name.
func (r *ParseResult) GetRulePosition(ruleName string) (diagnostic.Position, bool) {
	for _, rule := range r.Rules {
		if rule.Name == ruleName {
			return rule.Position, true
		}
	}
	return diagnostic.Position{}, false
}

// GetRuleInfo returns the RuleInfo for a rule by name.
func (r *ParseResult) GetRuleInfo(ruleName string) (*RuleInfo, bool) {
	for i := range r.Rules {
		if r.Rules[i].Name == ruleName {
			return &r.Rules[i], true
		}
	}
	return nil, false
}
