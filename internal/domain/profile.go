package domain

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadProfile reads and validates a profile YAML file. Validation covers the
// fields the rubric depends on: role targets, work constraints, and skill
// lists must be present and sane.
func LoadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read profile %s: %w", path, err)
	}
	return ParseProfile(data)
}

// ParseProfile unmarshals and validates profile YAML content.
func ParseProfile(data []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

// Validate enforces the profile contract used by scoring.
func (p *Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("profile: name is required")
	}
	if len(p.RoleTargets) == 0 {
		return fmt.Errorf("profile: at least one role_target is required")
	}
	if len(p.MustHaveSkills) == 0 {
		return fmt.Errorf("profile: at least one must_have_skill is required")
	}
	if len(p.WorkConstraints.Locations) == 0 && !p.WorkConstraints.RemoteOK {
		return fmt.Errorf("profile: work_constraints must have locations or remote_ok=true")
	}
	if p.WorkConstraints.MinYears < 0 {
		return fmt.Errorf("profile: min_years cannot be negative")
	}
	if len(p.WorkConstraints.WorkAuthorization) == 0 {
		return fmt.Errorf("profile: work_authorization is required")
	}
	if len(p.WorkConstraints.Languages) == 0 {
		return fmt.Errorf("profile: work_constraints.languages is required")
	}
	return nil
}
