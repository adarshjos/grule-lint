package parser

import (
	"github.com/hyperjumptech/grule-rule-engine/antlr/parser/grulev3"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

// RuleInfo contains parsed information about a single GRL rule.
type RuleInfo struct {
	Name        string
	Description string
	Salience    string
	Position    diagnostic.Position
	EndPosition diagnostic.Position

	WhenPosition diagnostic.Position
	ThenPosition diagnostic.Position

	FunctionCalls       []FunctionCallInfo
	ConditionCount      int
	HasWhenExpression   bool
	ThenActionCount     int
	VariableAssignments []VariableInfo
	VariableUsages      []VariableInfo
	WhenExpressionText  string
}

type FunctionCallInfo struct {
	Name     string
	Position diagnostic.Position
}

type VariableInfo struct {
	Name     string
	Position diagnostic.Position
}

// LintListener extracts rule information from the ANTLR parse tree.
type LintListener struct {
	*grulev3.Basegrulev3Listener
	Rules       []RuleInfo
	currentRule *RuleInfo
	inThenScope bool
	inWhenScope bool
}

func NewLintListener() *LintListener {
	return &LintListener{
		Basegrulev3Listener: &grulev3.Basegrulev3Listener{},
		Rules:               make([]RuleInfo, 0),
	}
}

func (l *LintListener) EnterRuleEntry(ctx *grulev3.RuleEntryContext) {
	l.currentRule = &RuleInfo{
		Position:      diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1},
		FunctionCalls: make([]FunctionCallInfo, 0),
	}
}

func (l *LintListener) ExitRuleEntry(ctx *grulev3.RuleEntryContext) {
	if l.currentRule != nil {
		l.currentRule.EndPosition = diagnostic.Position{Line: ctx.GetStop().GetLine(), Column: ctx.GetStop().GetColumn() + 1}
		l.Rules = append(l.Rules, *l.currentRule)
		l.currentRule = nil
	}
}

func (l *LintListener) EnterRuleName(ctx *grulev3.RuleNameContext) {
	if l.currentRule != nil {
		l.currentRule.Name = ctx.GetText()
	}
}

func (l *LintListener) EnterRuleDescription(ctx *grulev3.RuleDescriptionContext) {
	if l.currentRule != nil {
		text := ctx.GetText()
		if len(text) >= 2 {
			text = text[1 : len(text)-1]
		}
		l.currentRule.Description = text
	}
}

func (l *LintListener) EnterSalience(ctx *grulev3.SalienceContext) {
	if l.currentRule != nil {
		l.currentRule.Salience = ctx.GetText()
	}
}

func (l *LintListener) EnterWhenScope(ctx *grulev3.WhenScopeContext) {
	l.inWhenScope = true
	if l.currentRule != nil {
		l.currentRule.WhenPosition = diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1}
	}
}

func (l *LintListener) ExitWhenScope(ctx *grulev3.WhenScopeContext) {
	l.inWhenScope = false
}

// EnterComparisonOperator counts comparison operators (==, !=, <, >, <=, >=) in when clause.
// Only comparison operators are counted as conditions, not logical operators (&&, ||).
func (l *LintListener) EnterComparisonOperator(ctx *grulev3.ComparisonOperatorContext) {
	if l.currentRule != nil && l.inWhenScope {
		l.currentRule.ConditionCount++
	}
}

func (l *LintListener) EnterThenScope(ctx *grulev3.ThenScopeContext) {
	l.inThenScope = true
	if l.currentRule != nil {
		l.currentRule.ThenPosition = diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1}
	}
}

func (l *LintListener) ExitThenScope(ctx *grulev3.ThenScopeContext) {
	l.inThenScope = false
}

func (l *LintListener) EnterFunctionCall(ctx *grulev3.FunctionCallContext) {
	if l.currentRule != nil && l.inThenScope {
		l.currentRule.FunctionCalls = append(l.currentRule.FunctionCalls, FunctionCallInfo{
			Name:     ctx.SIMPLENAME().GetText(),
			Position: diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1},
		})
	}
}

func (l *LintListener) EnterExpression(ctx *grulev3.ExpressionContext) {
	if l.currentRule != nil && l.inWhenScope {
		l.currentRule.HasWhenExpression = true
		if l.currentRule.WhenExpressionText == "" {
			l.currentRule.WhenExpressionText = ctx.GetText()
		}
	}
}

// EnterAssignment tracks variable assignments in then clause for unused variable detection.
func (l *LintListener) EnterAssignment(ctx *grulev3.AssignmentContext) {
	if l.currentRule != nil && l.inThenScope {
		l.currentRule.ThenActionCount++
		// Track the assigned variable for GRL008 (unused-variable) detection
		if ctx.Variable() != nil {
			l.currentRule.VariableAssignments = append(l.currentRule.VariableAssignments, VariableInfo{
				Name:     ctx.Variable().GetText(),
				Position: diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1},
			})
		}
	}
}

func (l *LintListener) EnterVariable(ctx *grulev3.VariableContext) {
	if l.currentRule == nil {
		return
	}
	l.currentRule.VariableUsages = append(l.currentRule.VariableUsages, VariableInfo{
		Name:     ctx.GetText(),
		Position: diagnostic.Position{Line: ctx.GetStart().GetLine(), Column: ctx.GetStart().GetColumn() + 1},
	})
}
