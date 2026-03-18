package rules

import (
	"fmt"
	"strings"

	"github.com/karlhill/pipeguard/internal/parser"
)

type RequireStepRule struct {
	StepName string
}

func (r *RequireStepRule) ID() string          { return "require-step" }
func (r *RequireStepRule) Description() string { return "Required step not found: " + r.StepName }
func (r *RequireStepRule) Validate(config *parser.PipelineConfig) []Issue {
	found := false
	checkSteps := func(steps []parser.StepDefinition) {
		for _, sd := range steps {
			if sd.Step != nil && sd.Step.Name == r.StepName {
				found = true
			}
			if len(sd.Parallel) > 0 {
				for _, psd := range sd.Parallel {
					if psd.Step != nil && psd.Step.Name == r.StepName {
						found = true
					}
				}
			}
		}
	}

	checkSteps(config.Pipelines.Default)
	for _, steps := range config.Pipelines.Branches { checkSteps(steps) }
	for _, steps := range config.Pipelines.Custom { checkSteps(steps) }
	for _, steps := range config.Pipelines.PullRequests { checkSteps(steps) }
	for _, steps := range config.Pipelines.Tags { checkSteps(steps) }

	if !found {
		return []Issue{{
			RuleID:   r.ID(),
			Message:  "Required step '" + r.StepName + "' is missing from the pipeline configuration.",
			Severity: SeverityError,
		}}
	}
	return nil
}

type ForbidImageTagRule struct {
	ForbiddenTags []string
}

func (r *ForbidImageTagRule) ID() string          { return "forbid-image-tag" }
func (r *ForbidImageTagRule) Description() string { return "Forbidden image tag used" }
func (r *ForbidImageTagRule) Validate(config *parser.PipelineConfig) []Issue {
	var issues []Issue
	checkImage := func(image *parser.ImageDefinition, location string) {
		if image == nil { return }
		for _, tag := range r.ForbiddenTags {
			if strings.HasSuffix(image.Name, ":"+tag) || (!strings.Contains(image.Name, ":") && tag == "latest") {
				issues = append(issues, Issue{
					RuleID:   r.ID(),
					Message:  "Forbidden image tag '" + tag + "' used in " + location + ": " + image.Name,
					Severity: SeverityError,
				})
			}
		}
	}

	checkImage(config.Image, "global configuration")

	checkSteps := func(steps []parser.StepDefinition, loc string) {
		for _, sd := range steps {
			if sd.Step != nil {
				checkImage(sd.Step.Image, loc + " step '" + sd.Step.Name + "'")
			}
			for _, psd := range sd.Parallel {
				if psd.Step != nil {
					checkImage(psd.Step.Image, loc + " parallel step '" + psd.Step.Name + "'")
				}
			}
		}
	}

	checkSteps(config.Pipelines.Default, "default")
	for branch, steps := range config.Pipelines.Branches { checkSteps(steps, "branch '"+branch+"'") }
	for custom, steps := range config.Pipelines.Custom { checkSteps(steps, "custom '"+custom+"'") }

	return issues
}

type RequireManualTriggerRule struct {
	Deployment string
}

func (r *RequireManualTriggerRule) ID() string          { return "require-manual-trigger" }
func (r *RequireManualTriggerRule) Description() string { return "Manual trigger required for deployment: " + r.Deployment }
func (r *RequireManualTriggerRule) Validate(config *parser.PipelineConfig) []Issue {
	var issues []Issue
	checkSteps := func(steps []parser.StepDefinition, loc string) {
		for _, sd := range steps {
			if sd.Step != nil && sd.Step.Deployment == r.Deployment && sd.Step.Trigger != "manual" {
				issues = append(issues, Issue{
					RuleID:   r.ID(),
					Message:  "Deployment '" + r.Deployment + "' in " + loc + " must have 'trigger: manual'.",
					Severity: SeverityError,
				})
			}
		}
	}

	checkSteps(config.Pipelines.Default, "default")
	for branch, steps := range config.Pipelines.Branches { checkSteps(steps, "branch '"+branch+"'") }
	for custom, steps := range config.Pipelines.Custom { checkSteps(steps, "custom '"+custom+"'") }

	return issues
}

type AllowPipeListRule struct {
	AllowedPipes []string
}

func (r *AllowPipeListRule) ID() string          { return "allow-pipe-list" }
func (r *AllowPipeListRule) Description() string { return "Pipe usage restricted" }
func (r *AllowPipeListRule) Validate(config *parser.PipelineConfig) []Issue {
	var issues []Issue
	checkSteps := func(steps []parser.StepDefinition, loc string) {
		for _, sd := range steps {
			if sd.Step != nil {
				for _, entry := range sd.Step.Script {
					var pipeStr string
					switch v := entry.(type) {
					case string:
						pipeStr = v
					case map[string]interface{}:
						if p, ok := v["pipe"]; ok {
							pipeStr = fmt.Sprintf("%v", p)
						}
					}

					if strings.HasPrefix(pipeStr, "docker-public.packages.atlassian.com") || strings.Contains(pipeStr, "pipe:") {
						// Clean up pipe name if it's the full string or starts with pipe:
						pipeName := strings.TrimPrefix(pipeStr, "pipe:")
						
						allowed := false
						for _, allowedPipe := range r.AllowedPipes {
							if strings.HasPrefix(pipeName, allowedPipe) {
								allowed = true
								break
							}
						}
						if !allowed {
							issues = append(issues, Issue{
								RuleID:   r.ID(),
								Message:  "Forbidden pipe used in " + loc + ": " + pipeName,
								Severity: SeverityWarning,
							})
						}
					}
				}
			}
		}
	}

	checkSteps(config.Pipelines.Default, "default")
	for branch, steps := range config.Pipelines.Branches { checkSteps(steps, "branch '"+branch+"'") }
	for custom, steps := range config.Pipelines.Custom { checkSteps(steps, "custom '"+custom+"'") }

	return issues
}
