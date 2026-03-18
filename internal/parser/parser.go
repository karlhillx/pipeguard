package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type PipelineConfig struct {
	Image     *ImageDefinition       `yaml:"image,omitempty"`
	Pipelines PipelineDefinitions    `yaml:"pipelines"`
}

type ImageDefinition struct {
	Name string `yaml:"name"`
}

func (i *ImageDefinition) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		i.Name = value.Value
		return nil
	}
	type alias struct {
		Name string `yaml:"name"`
	}
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}
	i.Name = aux.Name
	return nil
}

type PipelineDefinitions struct {
	Default      []StepDefinition            `yaml:"default,omitempty"`
	Branches     map[string][]StepDefinition `yaml:"branches,omitempty"`
	Custom       map[string][]StepDefinition `yaml:"custom,omitempty"`
	PullRequests map[string][]StepDefinition `yaml:"pull-requests,omitempty"`
	Tags         map[string][]StepDefinition `yaml:"tags,omitempty"`
}

type StepDefinition struct {
	Step     *StepInner       `yaml:"step,omitempty"`
	Parallel []StepDefinition `yaml:"parallel,omitempty"`
}

type StepInner struct {
	Name       string           `yaml:"name,omitempty"`
	Image      *ImageDefinition `yaml:"image,omitempty"`
	Deployment string           `yaml:"deployment,omitempty"`
	Trigger    string           `yaml:"trigger,omitempty"`
	Script     []interface{}    `yaml:"script,omitempty"` // Use interface to handle both strings and maps (pipes)
}

func Parse(path string) (*PipelineConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config PipelineConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
