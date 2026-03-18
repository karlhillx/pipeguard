package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/karlhill/pipeguard/internal/parser"
	"github.com/karlhill/pipeguard/internal/report"
	"github.com/karlhill/pipeguard/internal/rules"
)

func main() {
	configPath := flag.String("config", "bitbucket-pipelines.yml", "Path to bitbucket-pipelines.yml")
	format := flag.String("format", "text", "Output format (text, json)")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s not found\n", *configPath)
		os.Exit(1)
	}

	pipeline, err := parser.Parse(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	engine := rules.NewEngine()
	// Add default MVP rules
	engine.AddRule(&rules.RequireStepRule{StepName: "test"})
	engine.AddRule(&rules.ForbidImageTagRule{ForbiddenTags: []string{"latest"}})
	engine.AddRule(&rules.RequireManualTriggerRule{Deployment: "production"})
	engine.AddRule(&rules.AllowPipeListRule{AllowedPipes: []string{"atlassian/slack-notify", "sonarsource/sonarcloud-scan"}})

	allIssues := engine.Run(pipeline)

	var formatter report.Formatter
	switch *format {
	case "json":
		formatter = &report.JSONFormatter{}
	default:
		formatter = &report.TextFormatter{}
	}

	if err := formatter.Format(allIssues, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}

	hasErrors := false
	for _, issue := range allIssues {
		if issue.Severity == rules.SeverityError {
			hasErrors = true
			break
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}
