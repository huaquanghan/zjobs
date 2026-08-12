package reporting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schema.json
var reportSchemaJSON string

var reportSchema = func() *jsonschema.Schema {
	c := jsonschema.NewCompiler()
	if err := c.AddResource("https://zjobs.local/report.schema.json",
		strings.NewReader(reportSchemaJSON)); err != nil {
		panic(fmt.Sprintf("reporting: compile schema: %v", err))
	}
	s, err := c.Compile("https://zjobs.local/report.schema.json")
	if err != nil {
		panic(fmt.Sprintf("reporting: compile schema: %v", err))
	}
	return s
}()

// ValidateJSON enforces the report output contract on rendered JSON.
func ValidateJSON(payload []byte) error {
	var raw any
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(&raw); err != nil {
		return fmt.Errorf("report: decode: %w", err)
	}
	if err := reportSchema.Validate(raw); err != nil {
		return fmt.Errorf("report: schema: %w", err)
	}
	return nil
}
