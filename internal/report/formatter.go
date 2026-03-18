package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/karlhill/pipeguard/internal/rules"
)

type Formatter interface {
	Format(issues []rules.Issue, w io.Writer) error
}

type TextFormatter struct{}

func (f *TextFormatter) Format(issues []rules.Issue, w io.Writer) error {
	if len(issues) == 0 {
		fmt.Fprintln(w, "✅ No policy violations found.")
		return nil
	}

	for _, issue := range issues {
		fmt.Fprintf(w, "[%-7s] (%s) %s\n", issue.Severity, issue.RuleID, issue.Message)
	}
	return nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(issues []rules.Issue, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(issues)
}
