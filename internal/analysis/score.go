package analysis

import (
	"fmt"
	"sort"
)

// Weights mirrors the rubric weights from config so analysis stays decoupled
// from the config package. Keys are the locked axis names.
type Weights map[string]float64

// Evaluate computes the deterministic verdict: any failed gate blocks to
// Skip; otherwise the weighted axis score maps to Apply / Stretch / Skip.
func Evaluate(a *Analysis, weights Weights) (VerdictResult, error) {
	// Hard gates: any explicit failure blocks Apply (R3).
	var failed []GateResult
	for _, g := range a.Gates {
		if !g.Passed {
			failed = append(failed, g)
		}
	}
	if len(failed) > 0 {
		sort.Slice(failed, func(i, j int) bool { return failed[i].Name < failed[j].Name })
		return VerdictResult{Verdict: VerdictSkip, Score: 0, FailedGates: failed}, nil
	}

	// Weighted 6-axis score over the configured weights.
	var weighted, weightSum float64
	for _, axis := range AxisNames {
		w, ok := weights[axis]
		if !ok || w < 0 {
			return VerdictResult{}, fmt.Errorf("analysis: missing or invalid weight for axis %q", axis)
		}
		fit, ok := a.Axes[axis]
		if !ok {
			return VerdictResult{}, fmt.Errorf("analysis: axis %q missing from payload", axis)
		}
		weighted += fit * w
		weightSum += w
	}
	if weightSum == 0 {
		return VerdictResult{}, fmt.Errorf("analysis: total rubric weight is zero")
	}
	score := weighted / weightSum

	verdict := VerdictSkip
	switch {
	case score >= applyThreshold:
		verdict = VerdictApply
	case score >= stretchThreshold:
		verdict = VerdictStretch
	}
	return VerdictResult{Verdict: verdict, Score: score}, nil
}
