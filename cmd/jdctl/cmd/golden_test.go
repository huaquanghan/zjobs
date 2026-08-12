package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"zjobs/internal/reporting"
)

// TestGoldenPipeline runs the locked CLI end-to-end against golden fixtures
// and compares the report byte-for-byte (JSON normalized only on
// generated_at, which is pinned to run time by design).
func TestGoldenPipeline(t *testing.T) {
	golden := filepath.Join("testdata", "golden")
	out := t.TempDir()

	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)
	rootCmd.SetArgs([]string{
		"analyze", "jd",
		"--profile", filepath.Join(golden, "profile-backend.yaml"),
		"--cv", filepath.Join(golden, "cv-main.md"),
		"--file", filepath.Join(golden, "jd-backend.txt"),
		"--analysis", filepath.Join(golden, "analysis.json"),
		"--out", out,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("analyze jd: %v", err)
	}
	if !strings.Contains(stdout.String(), "verdict=Apply") {
		t.Errorf("stdout = %q, want verdict=Apply", stdout.String())
	}

	entries, err := os.ReadDir(out)
	if err != nil {
		t.Fatalf("read out dir: %v", err)
	}
	var jsonPath, mdPath string
	for _, e := range entries {
		switch {
		case strings.HasSuffix(e.Name(), ".json"):
			jsonPath = filepath.Join(out, e.Name())
		case strings.HasSuffix(e.Name(), ".md"):
			mdPath = filepath.Join(out, e.Name())
		}
	}
	if jsonPath == "" || mdPath == "" {
		t.Fatalf("missing report files in %s", out)
	}

	// Report JSON must validate against the locked report schema.
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := reporting.ValidateJSON(raw); err != nil {
		t.Fatalf("report schema: %v", err)
	}

	// Byte-for-byte golden comparison, with the pinned timestamp normalized.
	var normalized struct {
		GeneratedAt any `json:"generated_at"`
	}
	if err := json.Unmarshal(raw, &normalized); err != nil {
		t.Fatal(err)
	}
	norm := bytes.Replace(raw, []byte(`"generated_at": `+string(mustJSON(t, normalized.GeneratedAt))), []byte(`"generated_at": "<pinned>"`), 1)
	compareGolden(t, filepath.Join(golden, "report.json"), norm)

	md, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	compareGolden(t, filepath.Join(golden, "report.md"), md)

	// Index must exist and contain the run.
	idx, err := os.ReadFile(filepath.Join(out, "index.jsonl"))
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	if len(idx) == 0 {
		t.Error("index.jsonl is empty")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func compareGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	want, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		t.Skipf("golden %s created (commit it); rerun to verify", path)
	}
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("golden mismatch for %s", path)
	}
}
