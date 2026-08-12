package domain

import (
	"path/filepath"
	"testing"
)

func TestLoadCVVariantValid(t *testing.T) {
	c, err := LoadCVVariant(filepath.Join("testdata", "cv-valid.md"))
	if err != nil {
		t.Fatalf("LoadCVVariant(valid): %v", err)
	}
	if c.Name != "cv-main" {
		t.Errorf("name = %q, want cv-main", c.Name)
	}
	if c.Profile != "backend" {
		t.Errorf("profile = %q, want backend", c.Profile)
	}
	if c.Body == "" || c.Body == " " {
		t.Error("body must be non-empty")
	}
}

func TestLoadCVVariantMissingFrontmatter(t *testing.T) {
	if _, err := LoadCVVariant(filepath.Join("testdata", "cv-invalid.md")); err == nil {
		t.Fatal("expected error for CV without frontmatter")
	}
}

func TestParseCVVariantEmptyBody(t *testing.T) {
	doc := "---\nname: cv\nprofile: p\n---\n   \n"
	if _, err := ParseCVVariant([]byte(doc)); err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestParseCVVariantMissingName(t *testing.T) {
	doc := "---\nprofile: p\n---\n# body\n"
	if _, err := ParseCVVariant([]byte(doc)); err == nil {
		t.Fatal("expected error for missing frontmatter name")
	}
}
