// Package config loads rubric weights and provider settings from YAML,
// merging over built-in defaults.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RubricWeights maps the six scoring axes to weights used by the analysis
// phase. Keys match the locked axis names.
type RubricWeights struct {
	HardConstraints  float64 `yaml:"hard_constraints" json:"hard_constraints"`
	MustHaveSkills   float64 `yaml:"must_have_skills" json:"must_have_skills"`
	NiceToHaveSkills float64 `yaml:"nice_to_have_skills" json:"nice_to_have_skills"`
	SeniorityScope   float64 `yaml:"seniority_scope" json:"seniority_scope"`
	DomainContext    float64 `yaml:"domain_context" json:"domain_context"`
	EvidenceStrength float64 `yaml:"evidence_strength" json:"evidence_strength"`
}

// ProviderSettings holds search provider configuration. Only Exa runs in the
// MVP; the key may come from the config file or the EXA_API_KEY env var.
type ProviderSettings struct {
	Exa struct {
		APIKey string `yaml:"api_key" json:"api_key"`
	} `yaml:"exa" json:"exa"`
}

// Config is the full runtime configuration.
type Config struct {
	Rubric   RubricWeights    `yaml:"rubric" json:"rubric"`
	Provider ProviderSettings `yaml:"provider" json:"provider"`
}

// Defaults returns the built-in configuration used when no file is given.
func Defaults() Config {
	return Config{
		Rubric: RubricWeights{
			HardConstraints:  1.0,
			MustHaveSkills:   1.0,
			NiceToHaveSkills: 0.5,
			SeniorityScope:   0.8,
			DomainContext:    0.6,
			EvidenceStrength: 0.7,
		},
	}
}

// Load reads a YAML config file and overlays it on Defaults. A missing file
// is not an error: defaults apply. An unreadable or malformed file is.
func Load(path string) (Config, error) {
	cfg := Defaults()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	return cfg, nil
}
