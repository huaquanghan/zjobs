// Package reporting renders stable Markdown + JSON reports and the
// append-only run index; it never writes to profile/CV source files.
package reporting

import (
	"time"

	"zjobs/internal/analysis"
	"zjobs/internal/domain"
)

// SchemaVersion is part of the stable report contract; bump only with a
// locked plan change.
const SchemaVersion = "1.0"

// Report is the single stable output schema shared by Markdown and JSON (R2).
// Field order is fixed by declaration order in JSON output.
type Report struct {
	SchemaVersion string                `json:"schema_version"`
	GeneratedAt   time.Time             `json:"generated_at"`
	Job           JobSummary            `json:"job"`
	Profile       string                `json:"profile"`
	CV            string                `json:"cv"`
	Verdict       string                `json:"verdict"`
	Score         float64               `json:"score"`
	Axes          map[string]float64    `json:"axes"`
	Gates         []analysis.GateResult `json:"gates"`
	Gaps          []analysis.Gap        `json:"gaps"`
	Prep          analysis.PrepPack     `json:"prep"`
}

// JobSummary is the normalized JD identity portion of a report.
type JobSummary struct {
	Title    string `json:"title,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	Source   string `json:"source"`
	URL      string `json:"url,omitempty"`
	Hash     string `json:"hash"`
}

// Build assembles a Report from the JD, profile/CV names, validated analysis,
// and the deterministic verdict.
func Build(jd *domain.JobDescription, profileName, cvName string, a *analysis.Analysis, res analysis.VerdictResult) *Report {
	return &Report{
		SchemaVersion: SchemaVersion,
		GeneratedAt:   time.Now().UTC(),
		Job: JobSummary{
			Title:    jd.Title,
			Company:  jd.Company,
			Location: jd.Location,
			Source:   string(jd.Source),
			URL:      jd.URL,
			Hash:     jd.Hash,
		},
		Profile: profileName,
		CV:      cvName,
		Verdict: res.Verdict,
		Score:   res.Score,
		Axes:    a.Axes,
		Gates:   a.Gates,
		Gaps:    a.Gaps,
		Prep:    a.Prep,
	}
}
