package analysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// schemaLoader compiles schema.json once per process.
var schemaLoader = func() *jsonschema.Schema {
	c := jsonschema.NewCompiler()
	if err := c.AddResource("https://zjobs.local/analysis.schema.json",
		strings.NewReader(schemaJSON)); err != nil {
		panic(fmt.Sprintf("analysis: compile schema: %v", err))
	}
	s, err := c.Compile("https://zjobs.local/analysis.schema.json")
	if err != nil {
		panic(fmt.Sprintf("analysis: compile schema: %v", err))
	}
	return s
}()

// Validate enforces the analysis contract: strict decoding (no unknown
// fields) plus JSON-schema checks on shape, ranges, and evidence presence.
func Validate(payload []byte) (*Analysis, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	var a Analysis
	if err := dec.Decode(&a); err != nil {
		return nil, fmt.Errorf("analysis: decode: %w", err)
	}
	var raw any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("analysis: decode: %w", err)
	}
	if err := schemaLoader.Validate(raw); err != nil {
		return nil, fmt.Errorf("analysis: schema: %w", err)
	}
	return &a, nil
}
