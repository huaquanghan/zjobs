package reporting

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWriteReportNeverTouchesProfileOrCV proves the read-only contract (R8):
// running a full write with real profile/CV fixture paths next to the out
// dir must leave those files byte-identical.
func TestWriteReportNeverTouchesProfileOrCV(t *testing.T) {
	work := t.TempDir()
	profilePath := filepath.Join(work, "profile.yaml")
	cvPath := filepath.Join(work, "cv.md")
	profileBefore := []byte("name: Tinh\nrole_targets: [x]\n")
	cvBefore := []byte("---\nname: cv-main\nprofile: p\n---\nbody\n")
	if err := os.WriteFile(profilePath, profileBefore, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cvPath, cvBefore, 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(work, "out")
	_, _, err := WriteReport(out, fixtureReport(t))
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}

	profileAfter, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read profile after: %v", err)
	}
	cvAfter, err := os.ReadFile(cvPath)
	if err != nil {
		t.Fatalf("read cv after: %v", err)
	}
	if string(profileAfter) != string(profileBefore) {
		t.Error("profile file was modified — violates R8")
	}
	if string(cvAfter) != string(cvBefore) {
		t.Error("cv file was modified — violates R8")
	}
}
