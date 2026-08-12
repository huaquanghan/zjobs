package reporting

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// IndexRow is one line of the append-only JSONL run index (R9).
type IndexRow struct {
	Timestamp time.Time `json:"ts"`
	JobHash   string    `json:"job_hash"`
	URL       string    `json:"url,omitempty"`
	Verdict   string    `json:"verdict"`
	Score     float64   `json:"score"`
}

// AppendIndex appends a run row to index.jsonl inside dir, skipping rows
// whose job hash already exists (dedupe on job hash).
func AppendIndex(dir string, row IndexRow) error {
	path := filepath.Join(dir, "index.jsonl")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("report: mkdir: %w", err)
	}
	hasHash, err := indexHasHash(path, row.JobHash)
	if err != nil {
		return err
	}
	if hasHash {
		return nil // dedupe: hash already recorded
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("report: open index: %w", err)
	}
	defer f.Close()
	line, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("report: encode index row: %w", err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("report: append index: %w", err)
	}
	return nil
}

func indexHasHash(path, hash string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("report: open index: %w", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 64*1024)
	for sc.Scan() {
		var row IndexRow
		if err := json.Unmarshal(sc.Bytes(), &row); err != nil {
			continue // skip corrupt lines rather than failing the run
		}
		if row.JobHash == hash {
			return true, nil
		}
	}
	return false, sc.Err()
}
