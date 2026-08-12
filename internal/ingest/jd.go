package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"zjobs/internal/domain"
)

// userAgent identifies the tool politely; pages that require more (login,
// CAPTCHA) are skipped by design (NG3).
const userAgent = "jdctl/0.1 (local-first job analysis; contact user)"

// FromFile reads a JD document from disk.
func FromFile(path string) (*domain.JobDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read jd file: %w", err)
	}
	return &domain.JobDescription{
		Body:      string(data),
		Source:    domain.SourceFile,
		FetchedAt: time.Now().UTC(),
		Hash:      hashJD(string(data), ""),
	}, nil
}

// FromPaste builds a JD from pasted text.
func FromPaste(text string) *domain.JobDescription {
	return &domain.JobDescription{
		Body:      strings.TrimSpace(text),
		Source:    domain.SourcePaste,
		FetchedAt: time.Now().UTC(),
		Hash:      hashJD(text, ""),
	}
}

// FromURL fetches a public JD page with a timeout, honoring robots.txt for
// the target host when it is reachable. No login or CAPTCHA bypass.
func FromURL(rawURL string, client *http.Client) (*domain.JobDescription, error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("jd url: %w", err)
	}
	if !u.IsAbs() || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("jd url: only absolute http(s) URLs are supported")
	}

	ok, err := robotsAllows(client, u)
	if err != nil {
		// A robots.txt outage must not silently gate a public page; the
		// respect signal is best-effort by design.
		ok = true
	}
	if !ok {
		return nil, fmt.Errorf("jd url: %s disallowed by robots.txt", u.String())
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jd fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jd fetch: %s returned %s", u.String(), resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf("jd read: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("jd fetch: empty body from %s", u.String())
	}
	return &domain.JobDescription{
		Body:      string(body),
		Source:    domain.SourceURL,
		URL:       u.String(),
		FetchedAt: time.Now().UTC(),
		Hash:      hashJD(string(body), u.String()),
	}, nil
}

// robotsAllows returns whether the URL's path is fetchable under the host's
// robots.txt (user-agent *). Missing robots.txt allows.
func robotsAllows(client *http.Client, u *url.URL) (bool, error) {
	robotsURL := *u
	robotsURL.Path = "/robots.txt"
	robotsURL.RawQuery = ""
	req, err := http.NewRequest(http.MethodGet, robotsURL.String(), nil)
	if err != nil {
		return true, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
		return true, nil
	}
	if resp.StatusCode != http.StatusOK {
		return true, nil
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return true, nil
	}
	return robotsAllowPath(string(data), u.Path), nil
}

// robotsAllowPath applies the simplest standard: a Disallow line whose value
// is a path prefix of targetPath blocks the fetch. Crawl-delay and
// per-agent groups beyond "*" are out of MVP scope.
func robotsAllowPath(robotsTxt, targetPath string) bool {
	if strings.TrimSpace(robotsTxt) == "" {
		return true
	}
	inStarGroup := false
	for _, line := range strings.Split(robotsTxt, "\n") {
		line = strings.TrimSpace(strings.ToLower(line))
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
			continue
		case strings.HasPrefix(line, "user-agent:"):
			inStarGroup = strings.TrimSpace(strings.TrimPrefix(line, "user-agent:")) == "*"
		case inStarGroup && strings.HasPrefix(line, "disallow:"):
			path := strings.TrimSpace(strings.TrimPrefix(line, "disallow:"))
			if path == "" {
				continue // empty Disallow means allow
			}
			if path == "/" || strings.HasPrefix(targetPath, path) {
				return false
			}
		}
	}
	return true
}

// hashJD gives a stable dedupe key: sha256 of source URL + normalized body.
func hashJD(body, sourceURL string) string {
	sum := sha256.Sum256([]byte(sourceURL + "\x00" + strings.ToLower(strings.TrimSpace(body))))
	return hex.EncodeToString(sum[:])
}
