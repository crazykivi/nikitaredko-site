package handlers

import (
	"testing"
)

func TestParseDateForSitemap_Valid(t *testing.T) {
	got := parseDateForSitemap("2025-06-15T14:00:00Z")
	if got != "2025-06-15" {
		t.Errorf("parseDateForSitemap = %q, want '2025-06-15'", got)
	}
}

func TestParseDateForSitemap_Invalid(t *testing.T) {
	got := parseDateForSitemap("garbage")
	// fallback = today's date, just check format YYYY-MM-DD
	if len(got) != 10 || got[4] != '-' || got[7] != '-' {
		t.Errorf("fallback date format invalid: %q", got)
	}
}
