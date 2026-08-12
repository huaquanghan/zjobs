package domain

import "time"

// Profile captures one user's job-search intent and constraints. It is the
// source of truth for scoring: role targets, work constraints (hard gates),
// and must-have/nice-to-have skills feed the 6-axis rubric.
type Profile struct {
	Name             string          `json:"name" yaml:"name"`
	RoleTargets      []string        `json:"role_targets" yaml:"role_targets"`
	WorkConstraints  WorkConstraints `json:"work_constraints" yaml:"work_constraints"`
	MustHaveSkills   []string        `json:"must_have_skills" yaml:"must_have_skills"`
	NiceToHaveSkills []string        `json:"nice_to_have_skills" yaml:"nice_to_have_skills"`
}

// WorkConstraints are the hard-gate inputs: violating any of them blocks an
// Apply verdict (R3).
type WorkConstraints struct {
	Locations         []string `json:"locations" yaml:"locations"`
	RemoteOK          bool     `json:"remote_ok" yaml:"remote_ok"`
	WorkAuthorization []string `json:"work_authorization" yaml:"work_authorization"`
	Languages         []string `json:"languages" yaml:"languages"`
	MinYears          int      `json:"min_years" yaml:"min_years"`
	Seniority         string   `json:"seniority" yaml:"seniority"`
}

// CVVariant is one Markdown CV bound to a profile. Frontmatter carries
// metadata; the body is the raw CV text used for evidence extraction.
type CVVariant struct {
	Name    string `json:"name" yaml:"name"`
	Profile string `json:"profile" yaml:"profile"`
	Body    string `json:"body"`
}

// JDSource tells where a JobDescription came from; exactly one is set.
type JDSource string

const (
	SourceURL   JDSource = "url"
	SourceFile  JDSource = "file"
	SourcePaste JDSource = "paste"
)

// JobDescription is one normalized job posting, whatever its origin.
type JobDescription struct {
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Body      string    `json:"body"`
	Source    JDSource  `json:"source"`
	URL       string    `json:"url,omitempty"`
	FetchedAt time.Time `json:"fetched_at"`
	Hash      string    `json:"hash"`
}
