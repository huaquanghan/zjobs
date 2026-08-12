package ingest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ExaProvider searches job postings via the Exa web-search API. It is the
// only real Provider in the MVP; the key comes from EXA_API_KEY or an
// injected API key, and the HTTP client is injectable for tests.
type ExaProvider struct {
	APIKey string
	Base   string
	Client *http.Client
}

// NewExaProvider builds a provider from the environment and config value,
// preferring the env var (never log the key).
func NewExaProvider(cfgAPIKey string) *ExaProvider {
	key := os.Getenv("EXA_API_KEY")
	if key == "" {
		key = cfgAPIKey
	}
	return &ExaProvider{
		APIKey: key,
		Base:   "https://api.exa.ai",
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

type exaRequest struct {
	Query      string `json:"query"`
	NumResults int    `json:"numResults"`
	Type       string `json:"type"`
}

type exaResponse struct {
	Results []exaResult `json:"results"`
}

type exaResult struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	Published string `json:"publishedDate"`
}

// Search runs one query against Exa's /search endpoint.
func (p *ExaProvider) Search(ctx context.Context, query string, limit int) ([]JobHit, error) {
	if p.APIKey == "" {
		return nil, fmt.Errorf("exa: EXA_API_KEY is not set")
	}
	if limit <= 0 {
		limit = 10
	}
	payload, err := json.Marshal(exaRequest{Query: query, NumResults: limit, Type: "auto"})
	if err != nil {
		return nil, fmt.Errorf("exa: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Base+"/search", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("User-Agent", userAgent)

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exa: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exa: %s", resp.Status)
	}
	var parsed exaResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("exa: decode response: %w", err)
	}
	hits := make([]JobHit, 0, len(parsed.Results))
	for _, r := range parsed.Results {
		hits = append(hits, JobHit{
			Title:   r.Title,
			URL:     r.URL,
			Company: r.Author,
			Snippet: r.Text,
		})
	}
	return hits, nil
}
