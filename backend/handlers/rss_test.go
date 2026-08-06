package handlers

import (
	"testing"
	"time"
)

func TestParseDateForRSS_ValidDate(t *testing.T) {
	got := parseDateForRSS("2025-01-15T10:30:00Z")

	expected := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC).Format(time.RFC1123Z)
	if got != expected {
		t.Errorf("parseDateForRSS = %q, want %q", got, expected)
	}
}

func TestParseDateForRSS_InvalidDate(t *testing.T) {
	got := parseDateForRSS("not-a-date")
	if got == "" {
		t.Error("fallback should produce non-empty date")
	}
}

func TestParseDateForRSS_EmptyString(t *testing.T) {
	got := parseDateForRSS("")
	if got == "" {
		t.Error("empty input should fallback to current time")
	}
}
