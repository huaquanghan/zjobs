package reporting

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteReport persists the report into dir as report-{hash8}.json and
// report-{hash8}.md, and appends the run to index.jsonl. It takes only the
// out directory and report data — profile/CV paths never enter here, so the
// tool structurally cannot modify them (R8).
func WriteReport(dir string, r *Report) (jsonPath, mdPath string, err error) {
	prefix := "report-" + safePrefix(r.Job.Hash)
	jsonPath = filepath.Join(dir, prefix+".json")
	mdPath = filepath.Join(dir, prefix+".md")

	js, err := RenderJSON(r)
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", fmt.Errorf("report: mkdir: %w", err)
	}
	if err := os.WriteFile(jsonPath, js, 0o644); err != nil {
		return "", "", fmt.Errorf("report: write json: %w", err)
	}
	if err := os.WriteFile(mdPath, RenderMarkdown(r), 0o644); err != nil {
		return "", "", fmt.Errorf("report: write md: %w", err)
	}

	if err := AppendIndex(dir, IndexRow{
		Timestamp: r.GeneratedAt,
		JobHash:   r.Job.Hash,
		URL:       r.Job.URL,
		Verdict:   r.Verdict,
		Score:     r.Score,
	}); err != nil {
		return "", "", err
	}
	return jsonPath, mdPath, nil
}

func safePrefix(hash string) string {
	if len(hash) >= 8 {
		return hash[:8]
	}
	return hash
}
