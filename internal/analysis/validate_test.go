package analysis

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func TestValidateValidApply(t *testing.T) {
	a, err := Validate(loadFixture(t, "analysis-apply.json"))
	if err != nil {
		t.Fatalf("Validate(apply): %v", err)
	}
	if len(a.Gates) != 5 {
		t.Errorf("gates = %d, want 5", len(a.Gates))
	}
}

func TestValidateRejectsMissingEvidence(t *testing.T) {
	_, err := Validate(loadFixture(t, "analysis-malformed.json"))
	if err == nil {
		t.Fatal("expected error for gap without jd_evidence/cv_evidence")
	}
}

func TestValidateRejectsUnknownField(t *testing.T) {
	payload := strings.Replace(string(loadFixture(t, "analysis-apply.json")),
		`"prep"`, `"bogus": 1, "prep"`, 1)
	if _, err := Validate([]byte(payload)); err == nil {
		t.Fatal("expected error for unknown top-level field")
	}
}

func TestValidateRejectsBadAxis(t *testing.T) {
	payload := strings.Replace(string(loadFixture(t, "analysis-apply.json")),
		`"hard_constraints": 0.9`, `"hard_constraints": 1.5`, 1)
	if _, err := Validate([]byte(payload)); err == nil {
		t.Fatal("expected error for axis out of range")
	}
}

func TestValidateRejectsUnknownGateName(t *testing.T) {
	payload := strings.Replace(string(loadFixture(t, "analysis-apply.json")),
		`"location"`, `"teleport"`, 1)
	if _, err := Validate([]byte(payload)); err == nil {
		t.Fatal("expected error for unknown gate name")
	}
}

func TestValidateRejectsBadAction(t *testing.T) {
	payload := strings.Replace(string(loadFixture(t, "analysis-apply.json")),
		`"action": "learn"`, `"action": "wing-it"`, 1)
	if _, err := Validate([]byte(payload)); err == nil {
		t.Fatal("expected error for unknown gap action")
	}
}

func TestValidateRejectsUnverifiedReference(t *testing.T) {
	payload := strings.Replace(string(loadFixture(t, "analysis-apply.json")),
		`"url": "https://example.com/postgres-scale"`, `"url": "ftp://example.com/x"`, 1)
	if _, err := Validate([]byte(payload)); err == nil {
		t.Fatal("expected error for non-http reference URL")
	}
}
