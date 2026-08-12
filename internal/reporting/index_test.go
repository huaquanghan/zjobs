package reporting

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"zjobs/internal/analysis"
)

func TestAppendIndexDedupesOnHash(t *testing.T) {
	dir := t.TempDir()
	row := IndexRow{JobHash: "aaa111", Verdict: analysis.VerdictApply, Score: 0.9}
	if err := AppendIndex(dir, row); err != nil {
		t.Fatalf("append 1: %v", err)
	}
	if err := AppendIndex(dir, row); err != nil {
		t.Fatalf("append 2 (dup): %v", err)
	}
	other := row
	other.JobHash = "bbb222"
	if err := AppendIndex(dir, other); err != nil {
		t.Fatalf("append 3: %v", err)
	}
	lines := countLines(t, filepath.Join(dir, "index.jsonl"))
	if lines != 2 {
		t.Errorf("index lines = %d, want 2 (dedupe on hash)", lines)
	}
}

func TestAppendIndexCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "reports")
	if err := AppendIndex(dir, IndexRow{JobHash: "c3", Verdict: analysis.VerdictSkip}); err != nil {
		t.Fatalf("append: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "index.jsonl")); err != nil {
		t.Fatalf("index not created: %v", err)
	}
}

func countLines(t *testing.T, path string) int {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open index: %v", err)
	}
	defer f.Close()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		n++
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan index: %v", err)
	}
	return n
}
