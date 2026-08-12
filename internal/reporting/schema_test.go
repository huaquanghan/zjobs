package reporting

import (
	"strings"
	"testing"
)

func TestValidateJSONAcceptsRenderedReport(t *testing.T) {
	payload, err := RenderJSON(fixtureReport(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateJSON(payload); err != nil {
		t.Fatalf("ValidateJSON(rendered): %v", err)
	}
}

func TestValidateJSONRejectsBadVerdict(t *testing.T) {
	payload, _ := RenderJSON(fixtureReport(t))
	bad := strings.Replace(string(payload), `"verdict": "Apply"`, `"verdict": "Maybe"`, 1)
	if err := ValidateJSON([]byte(bad)); err == nil {
		t.Fatal("expected error for unknown verdict")
	}
}

func TestValidateJSONRejectsBadVersion(t *testing.T) {
	payload, _ := RenderJSON(fixtureReport(t))
	bad := strings.Replace(string(payload), `"schema_version": "1.0"`, `"schema_version": "9.9"`, 1)
	if err := ValidateJSON([]byte(bad)); err == nil {
		t.Fatal("expected error for unknown schema version")
	}
}

func TestValidateJSONRejectsUnknownField(t *testing.T) {
	payload, _ := RenderJSON(fixtureReport(t))
	bad := strings.Replace(string(payload), `"prep": {`, `"sneaky": 1, "prep": {`, 1)
	if err := ValidateJSON([]byte(bad)); err == nil {
		t.Fatal("expected error for unknown field")
	}
}
