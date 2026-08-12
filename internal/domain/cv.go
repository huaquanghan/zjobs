package domain

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadCVVariant reads a CV variant Markdown file with YAML frontmatter:
//
//	---
//	name: cv-main
//	profile: backend
//	---
//	# CV body...
func LoadCVVariant(path string) (*CVVariant, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read cv %s: %w", path, err)
	}
	return ParseCVVariant(data)
}

// ParseCVVariant splits frontmatter from body and validates both parts.
func ParseCVVariant(data []byte) (*CVVariant, error) {
	text := string(data)
	raw, body, ok := splitFrontmatter(text)
	if !ok {
		return nil, fmt.Errorf("cv: missing YAML frontmatter (--- blocks)")
	}
	var c CVVariant
	if err := yaml.Unmarshal([]byte(raw), &c); err != nil {
		return nil, fmt.Errorf("cv: parse frontmatter: %w", err)
	}
	c.Body = strings.TrimSpace(body)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate enforces the CV variant contract.
func (c *CVVariant) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("cv: frontmatter name is required")
	}
	if strings.TrimSpace(c.Profile) == "" {
		return fmt.Errorf("cv: frontmatter profile is required")
	}
	if strings.TrimSpace(c.Body) == "" {
		return fmt.Errorf("cv: body cannot be empty")
	}
	return nil
}

// splitFrontmatter returns (frontmatter, body) when the document starts with
// a --- line and closes it at the next --- line.
func splitFrontmatter(text string) (string, string, bool) {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return "", "", false
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			fm := strings.Join(lines[1:i], "\n")
			body := strings.Join(lines[i+1:], "\n")
			return fm, body, true
		}
	}
	return "", "", false
}
