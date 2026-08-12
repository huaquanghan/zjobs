package ingest

import (
	"context"
	"time"

	"zjobs/internal/domain"
)

// JobHit is one search result from a provider, before full JD ingestion.
type JobHit struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Company     string    `json:"company"`
	Snippet     string    `json:"snippet,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

// Provider discovers job hits for a query. Implementations must be safe to
// call from tests via a mockable HTTP client; the MVP ships only Exa (R6).
type Provider interface {
	Search(ctx context.Context, query string, limit int) ([]JobHit, error)
}

// NormalizeHit converts a provider hit into a minimal JobDescription shell;
// full JD content lands after the fetching phase.
func NormalizeHit(h JobHit) domain.JobDescription {
	return domain.JobDescription{
		Title:     h.Title,
		Company:   h.Company,
		Source:    domain.SourceURL,
		URL:       h.URL,
		FetchedAt: time.Now().UTC(),
	}
}
