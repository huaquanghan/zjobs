package ingest

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"zjobs/internal/domain"
)

func TestFromFile(t *testing.T) {
	jd, err := FromFile(filepath.Join("testdata", "jd.txt"))
	if err != nil {
		t.Fatalf("FromFile: %v", err)
	}
	if jd.Source != domain.SourceFile {
		t.Errorf("source = %q, want file", jd.Source)
	}
	if jd.Body == "" {
		t.Error("body must be non-empty")
	}
	if jd.Hash == "" {
		t.Error("hash must be set")
	}
}

func TestFromPaste(t *testing.T) {
	jd := FromPaste("  Senior Go engineer at Acme  ")
	if jd.Source != domain.SourcePaste {
		t.Errorf("source = %q, want paste", jd.Source)
	}
	if strings.TrimSpace(jd.Body) != jd.Body {
		t.Error("body must be trimmed")
	}
}

func TestFromURLHappyPath(t *testing.T) {
	body := "<html><body>Senior Go engineer at Acme</body></html>"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path == "/jobs/1" {
			_, _ = w.Write([]byte(body))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	jd, err := FromURL(srv.URL+"/jobs/1", srv.Client())
	if err != nil {
		t.Fatalf("FromURL: %v", err)
	}
	if jd.Source != domain.SourceURL {
		t.Errorf("source = %q, want url", jd.Source)
	}
	if jd.URL != srv.URL+"/jobs/1" {
		t.Errorf("url = %q", jd.URL)
	}
	if jd.Body != body {
		t.Errorf("body mismatch: got %q", jd.Body)
	}
}

func TestFromURLRobotsDisallowed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte("User-agent: *\nDisallow: /private/\n"))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	if _, err := FromURL(srv.URL+"/private/jobs/1", srv.Client()); err == nil {
		t.Fatal("expected error for robots-disallowed path")
	}
}

func TestFromURLNonHTTP(t *testing.T) {
	if _, err := FromURL("ftp://example.com/jobs/1", nil); err == nil {
		t.Fatal("expected error for non-http scheme")
	}
}

func TestFromURLMissingFile(t *testing.T) {
	if _, err := FromFile(filepath.Join("testdata", "nope.txt")); err == nil {
		t.Fatal("expected error for missing file")
	}
}
