package analysis

import (
	"testing"
)

var defaultWeights = Weights{
	AxisHardConstraints:  1.0,
	AxisMustHaveSkills:   1.0,
	AxisNiceToHaveSkills: 0.5,
	AxisSeniorityScope:   0.8,
	AxisDomainContext:    0.6,
	AxisEvidenceStrength: 0.7,
}

func TestEvaluateApply(t *testing.T) {
	a, err := Validate(loadFixture(t, "analysis-apply.json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Evaluate(a, defaultWeights)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if res.Verdict != VerdictApply {
		t.Errorf("verdict = %q, want Apply (score %.2f)", res.Verdict, res.Score)
	}
	if len(res.FailedGates) != 0 {
		t.Errorf("failed gates = %d, want 0", len(res.FailedGates))
	}
}

func TestEvaluateGateFailFlipsToSkip(t *testing.T) {
	a, err := Validate(loadFixture(t, "analysis-gate-fail.json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Evaluate(a, defaultWeights)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if res.Verdict != VerdictSkip {
		t.Errorf("verdict = %q, want Skip despite high axes", res.Verdict)
	}
	if len(res.FailedGates) != 1 || res.FailedGates[0].Name != GateLocation {
		t.Errorf("failed gates = %+v, want location only", res.FailedGates)
	}
}

func TestEvaluateStretch(t *testing.T) {
	a, err := Validate(loadFixture(t, "analysis-stretch.json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := Evaluate(a, defaultWeights)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if res.Verdict != VerdictStretch {
		t.Errorf("verdict = %q, want Stretch (score %.2f)", res.Verdict, res.Score)
	}
}

// TestEachGateCanFlipVerdict proves every gate name can independently block
// to Skip (R3), one gate at a time.
func TestEachGateCanFlipVerdict(t *testing.T) {
	base, err := Validate(loadFixture(t, "analysis-apply.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range GateNames {
		cp := *base
		gates := make([]GateResult, len(base.Gates))
		copy(gates, base.Gates)
		found := false
		for i := range gates {
			if gates[i].Name == name {
				gates[i].Passed = false
				gates[i].Reason = "test flip"
				gates[i].JDEvidence = "JD evidence for " + name
				found = true
			}
		}
		if !found {
			t.Fatalf("gate %q not present in fixture", name)
		}
		cp.Gates = gates
		res, err := Evaluate(&cp, defaultWeights)
		if err != nil {
			t.Fatalf("Evaluate(%s): %v", name, err)
		}
		if res.Verdict != VerdictSkip {
			t.Errorf("gate %q did not flip verdict: got %q", name, res.Verdict)
		}
	}
}

func TestEvaluateMissingWeight(t *testing.T) {
	a, err := Validate(loadFixture(t, "analysis-apply.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Evaluate(a, Weights{AxisMustHaveSkills: 1.0}); err == nil {
		t.Fatal("expected error for missing axis weight")
	}
}
