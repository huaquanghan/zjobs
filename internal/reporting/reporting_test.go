package reporting

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"zjobs/internal/analysis"
	"zjobs/internal/domain"
)

func fixtureAnalysis(t *testing.T) (*analysis.Analysis, analysis.VerdictResult) {
	t.Helper()
	payload, err := os.ReadFile(filepath.Join("..", "analysis", "testdata", "analysis-apply.json"))
	if err != nil {
		t.Fatalf("read analysis fixture: %v", err)
	}
	a, err := analysis.Validate(payload)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	res, err := analysis.Evaluate(a, analysis.Weights{
		analysis.AxisHardConstraints:  1.0,
		analysis.AxisMustHaveSkills:   1.0,
		analysis.AxisNiceToHaveSkills: 0.5,
		analysis.AxisSeniorityScope:   0.8,
		analysis.AxisDomainContext:    0.6,
		analysis.AxisEvidenceStrength: 0.7,
	})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	return a, res
}

func fixtureReport(t *testing.T) *Report {
	t.Helper()
	a, res := fixtureAnalysis(t)
	jd := &domain.JobDescription{
		Title:    "Senior Backend Engineer (Go)",
		Company:  "Acme",
		Location: "Hanoi, hybrid",
		Source:   domain.SourceURL,
		URL:      "https://acme.dev/jobs/1",
		Hash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	r := Build(jd, "backend", "cv-main", a, res)
	r.GeneratedAt = time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC) // stable golden
	return r
}

// TestRenderGoldenJSON proves the JSON renderer is stable byte-for-byte.
func TestRenderGoldenJSON(t *testing.T) {
	got, err := RenderJSON(fixtureReport(t))
	if err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	compareGolden(t, "report-apply.json", got)
}

// TestRenderGoldenMarkdown proves the Markdown renderer is stable.
func TestRenderGoldenMarkdown(t *testing.T) {
	got := RenderMarkdown(fixtureReport(t))
	compareGolden(t, "report-apply.md", got)
}

func compareGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	want, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// First run: materialize the golden so it is committed with the repo.
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden: %v", err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", name, err)
		}
		t.Skipf("golden %s created (commit it); rerun to verify", name)
	}
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("golden mismatch for %s", name)
	}
}
