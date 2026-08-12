package domain

import (
	"path/filepath"
	"testing"
)

func TestLoadProfileValid(t *testing.T) {
	p, err := LoadProfile(filepath.Join("testdata", "profile-valid.yaml"))
	if err != nil {
		t.Fatalf("LoadProfile(valid): %v", err)
	}
	if p.Name != "Tinh" {
		t.Errorf("name = %q, want Tinh", p.Name)
	}
	if len(p.MustHaveSkills) != 3 {
		t.Errorf("must_have_skills = %d entries, want 3", len(p.MustHaveSkills))
	}
	if !p.WorkConstraints.RemoteOK {
		t.Error("remote_ok = false, want true")
	}
	if p.WorkConstraints.MinYears != 5 {
		t.Errorf("min_years = %d, want 5", p.WorkConstraints.MinYears)
	}
}

func TestLoadProfileInvalidMissingSkills(t *testing.T) {
	_, err := LoadProfile(filepath.Join("testdata", "profile-invalid.yaml"))
	if err == nil {
		t.Fatal("LoadProfile(invalid): expected error, got nil")
	}
}

func TestParseProfileMissingName(t *testing.T) {
	_, err := ParseProfile([]byte("role_targets: [x]\nwork_constraints:\n  locations: [Hanoi]\n  work_authorization: [VN]\n  languages: [vi]\nmust_have_skills: [Go]\n"))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParseProfileNoRemoteNoLocation(t *testing.T) {
	_, err := ParseProfile([]byte("name: X\nrole_targets: [x]\nwork_constraints:\n  remote_ok: false\n  work_authorization: [VN]\n  languages: [vi]\nmust_have_skills: [Go]\n"))
	if err == nil {
		t.Fatal("expected error when no locations and remote_ok=false")
	}
}

func TestLoadProfileMissingFile(t *testing.T) {
	if _, err := LoadProfile(filepath.Join("testdata", "nope.yaml")); err == nil {
		t.Fatal("expected error for missing file")
	}
}
