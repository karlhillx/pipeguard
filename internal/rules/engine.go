package rules

import (
	"github.com/karlhill/pipeguard/internal/parser"
)

type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

type Issue struct {
	RuleID   string
	Message  string
	Severity Severity
}

type Rule interface {
	ID() string
	Description() string
	Validate(config *parser.PipelineConfig) []Issue
}

type Engine struct {
	rules []Rule
}

func NewEngine() *Engine {
	return &Engine{rules: []Rule{}}
}

func (e *Engine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

func (e *Engine) Run(config *parser.PipelineConfig) []Issue {
	var allIssues []Issue
	for _, rule := range e.rules {
		issues := rule.Validate(config)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}
