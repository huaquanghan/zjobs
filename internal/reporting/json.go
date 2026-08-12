package reporting

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// RenderJSON serializes the report with indentation and a trailing newline,
// in the fixed field order of the Report struct.
func RenderJSON(r *Report) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return nil, fmt.Errorf("report: encode json: %w", err)
	}
	return buf.Bytes(), nil
}
