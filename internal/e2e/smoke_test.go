// Package e2e proves the locked contract on a real public JD: fetch without
// login or CAPTCHA bypass, build a report, validate it against the schema,
// and confirm profile/CV files are untouched.
package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"zjobs/internal/analysis"
	"zjobs/internal/domain"
	"zjobs/internal/ingest"
	"zjobs/internal/reporting"
)

// smokeJDURL is a public job board reachable without login, robots-clean
// (only /embed/ is disallowed), served server-side HTML.
const smokeJDURL = "https://boards.greenhouse.io/gitlab"

var defaultWeights = analysis.Weights{
	analysis.AxisHardConstraints:  1.0,
	analysis.AxisMustHaveSkills:   1.0,
	analysis.AxisNiceToHaveSkills: 0.5,
	analysis.AxisSeniorityScope:   0.8,
	analysis.AxisDomainContext:    0.6,
	analysis.AxisEvidenceStrength: 0.7,
}

// TestSmokeRealJD is the live smoke run (R10). It needs network, so it runs
// only when RUN_JDSMOKE=1. The produced report is committed as a fixture;
// TestSmokeFixtureReportValidates re-validates it hermetically.
func TestSmokeRealJD(t *testing.T) {
	if os.Getenv("RUN_JDSMOKE") != "1" {
		t.Skip("set RUN_JDSMOKE=1 to run the live fetch smoke test")
	}

	work := t.TempDir()
	profileSrc := filepath.Join("..", "..", "cmd", "jdctl", "cmd", "testdata", "golden", "profile-backend.yaml")
	cvSrc := filepath.Join("..", "..", "cmd", "jdctl", "cmd", "testdata", "golden", "cv-main.md")
	profilePath := filepath.Join(work, "profile.yaml")
	cvPath := filepath.Join(work, "cv.md")
	for dst, src := range map[string]string{profilePath: profileSrc, cvPath: cvSrc} {
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	profileBefore, _ := os.ReadFile(profilePath)
	cvBefore, _ := os.ReadFile(cvPath)

	// Fetch the real public JD (public fetch only, robots respected).
	jd, err := ingest.FromURL(smokeJDURL, nil)
	if err != nil {
		t.Fatalf("fetch real JD: %v", err)
	}
	profile, err := domain.LoadProfile(profilePath)
	if err != nil {
		t.Fatal(err)
	}
	cv, err := domain.LoadCVVariant(cvPath)
	if err != nil {
		t.Fatal(err)
	}

	// Semantic layer comes from the Claude skill; for the smoke run it is a
	// committed golden analysis payload.
	payload, err := os.ReadFile(filepath.Join("..", "..", "cmd", "jdctl", "cmd", "testdata", "golden", "analysis.json"))
	if err != nil {
		t.Fatal(err)
	}
	ana, err := analysis.Validate(payload)
	if err != nil {
		t.Fatal(err)
	}
	res, err := analysis.Evaluate(ana, defaultWeights)
	if err != nil {
		t.Fatal(err)
	}

	report := reporting.Build(jd, profile.Name, cv.Name, ana, res)
	jsonPayload, err := reporting.RenderJSON(report)
	if err != nil {
		t.Fatal(err)
	}
	if err := reporting.ValidateJSON(jsonPayload); err != nil {
		t.Fatalf("smoke report fails schema: %v", err)
	}

	// Profile/CV must be byte-identical after the whole run (R8).
	profileAfter, _ := os.ReadFile(profilePath)
	cvAfter, _ := os.ReadFile(cvPath)
	if string(profileAfter) != string(profileBefore) {
		t.Error("profile modified by smoke run — violates R8")
	}
	if string(cvAfter) != string(cvBefore) {
		t.Error("cv modified by smoke run — violates R8")
	}

	// Commit the report as a fixture for hermetic re-validation.
	dir := filepath.Join("testdata", "smoke")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report.json"), jsonPayload, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report.md"), reporting.RenderMarkdown(report), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("smoke report committed: url=%s bytes=%d verdict=%s score=%.2f",
		jd.URL, len(jsonPayload), res.Verdict, res.Score)
}

// TestSmokeFixtureReportValidates re-validates the committed smoke report
// without network; it runs as part of the normal suite.
func TestSmokeFixtureReportValidates(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("testdata", "smoke", "report.json"))
	if err != nil {
		t.Fatalf("missing committed smoke fixture (run RUN_JDSMOKE=1 once): %v", err)
	}
	if err := reporting.ValidateJSON(payload); err != nil {
		t.Fatalf("committed smoke report fails schema: %v", err)
	}
}
