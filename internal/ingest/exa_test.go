package ingest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExaSearchUsesMockedClient proves the adapter works against the Exa
// response shape without touching the network.
func TestExaSearchUsesMockedClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("authorization = %q", got)
		}
		var req exaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if req.Query != "Go backend Hanoi" {
			t.Errorf("query = %q", req.Query)
		}
		_ = json.NewEncoder(w).Encode(exaResponse{
			Results: []exaResult{
				{Title: "Senior Go Engineer", URL: "https://acme.dev/jobs/1", Author: "Acme", Text: "Go, PostgreSQL"},
			},
		})
	}))
	defer srv.Close()

	p := &ExaProvider{APIKey: "test-key", Base: srv.URL, Client: srv.Client()}
	hits, err := p.Search(context.Background(), "Go backend Hanoi", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(hits) != 1 {
		t.Fatalf("hits = %d, want 1", len(hits))
	}
	if hits[0].Title != "Senior Go Engineer" || hits[0].Company != "Acme" {
		t.Errorf("hit = %+v", hits[0])
	}
}

func TestExaSearchMissingKey(t *testing.T) {
	p := &ExaProvider{Base: "http://example.invalid"}
	if _, err := p.Search(context.Background(), "q", 5); err == nil {
		t.Fatal("expected error when API key missing")
	}
}

func TestExaSearchErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTooManyRequests)
	}))
	defer srv.Close()
	p := &ExaProvider{APIKey: "k", Base: srv.URL, Client: srv.Client()}
	if _, err := p.Search(context.Background(), "q", 5); err == nil {
		t.Fatal("expected error on non-200")
	}
}
