package analysis

import _ "embed"

// schemaJSON is the locked analysis contract document, embedded so the
// validator and the Claude skill always share one source of truth.
//
//go:embed schema.json
var schemaJSON string
