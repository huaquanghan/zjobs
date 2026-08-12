// Package analysis implements the 6-axis rubric, hard gates, evidence-backed
// gaps, and the Apply/Stretch/Skip verdict over the Claude analysis contract.
package analysis

// Axis names of the locked rubric (R3). JSON keys match the config weights.
const (
	AxisHardConstraints  = "hard_constraints"
	AxisMustHaveSkills   = "must_have_skills"
	AxisNiceToHaveSkills = "nice_to_have_skills"
	AxisSeniorityScope   = "seniority_scope"
	AxisDomainContext    = "domain_context"
	AxisEvidenceStrength = "evidence_strength"
)

// AxisNames lists all six axes in a stable order.
var AxisNames = []string{
	AxisHardConstraints,
	AxisMustHaveSkills,
	AxisNiceToHaveSkills,
	AxisSeniorityScope,
	AxisDomainContext,
	AxisEvidenceStrength,
}

// Gate names for work-constraint hard gates (R3).
const (
	GateLocation          = "location"
	GateWorkAuthorization = "work_authorization"
	GateLanguage          = "language"
	GateExperienceYears   = "experience_years"
	GateSeniority         = "seniority"
)

// GateNames lists all gates in a stable order.
var GateNames = []string{
	GateLocation,
	GateWorkAuthorization,
	GateLanguage,
	GateExperienceYears,
	GateSeniority,
}

// GapAction is the normalized action attached to an evidence-backed gap (R4).
type GapAction string

const (
	ActionLearn        GapAction = "learn"
	ActionRewriteCV    GapAction = "rewrite_cv"
	ActionPrepareStory GapAction = "prepare_story"
	ActionSkip         GapAction = "skip"
)

// Analysis is the structured JSON contract the Claude skill emits. Go
// validates it against schema.json and computes the deterministic verdict.
type Analysis struct {
	Axes  map[string]float64 `json:"axes"`
	Gates []GateResult       `json:"gates"`
	Gaps  []Gap              `json:"gaps"`
	Prep  PrepPack           `json:"prep"`
}

// GateResult is one work-constraint gate check. When Passed is false the
// verdict is blocked to Skip; Reason and JDEvidence must still be filled.
type GateResult struct {
	Name       string `json:"name"`
	Passed     bool   `json:"passed"`
	Reason     string `json:"reason"`
	JDEvidence string `json:"jd_evidence"`
}

// Gap is one missing JD requirement with citations from both sides and an
// action. Evidence fields must be non-empty (R4).
type Gap struct {
	Requirement string    `json:"requirement"`
	JDEvidence  string    `json:"jd_evidence"`
	CVEvidence  string    `json:"cv_evidence"`
	Impact      string    `json:"impact"`
	Action      GapAction `json:"action"`
}

// PrepPack is the lite preparation pack (R5): 3–5 questions, verified
// references only, and 1–2 next actions.
type PrepPack struct {
	Questions   []string    `json:"questions"`
	References  []Reference `json:"references"`
	NextActions []string    `json:"next_actions"`
}

// Reference is a verified URL: title, URL, and the source that surfaced it.
// Never invented (R5).
type Reference struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

// VerdictResult is the deterministic outcome of Evaluate.
type VerdictResult struct {
	Verdict     string       `json:"verdict"`
	Score       float64      `json:"score"`
	FailedGates []GateResult `json:"failed_gates,omitempty"`
}

// Verdict labels (locked).
const (
	VerdictApply   = "Apply"
	VerdictStretch = "Stretch"
	VerdictSkip    = "Skip"
)

// Score thresholds: Apply >= applyThreshold, Stretch >= stretchThreshold.
const (
	applyThreshold   = 0.8
	stretchThreshold = 0.6
)
