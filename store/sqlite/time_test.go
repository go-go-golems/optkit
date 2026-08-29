package sqlite

import (
	"testing"
	"time"
)

func TestFormatTimeIsLexicallyChronologicalWithinSecond(t *testing.T) {
	whole := time.Date(2026, time.August, 29, 12, 0, 0, 0, time.UTC)
	fractional := whole.Add(500 * time.Millisecond)
	wholeText := formatTime(whole)
	fractionalText := formatTime(fractional)
	if wholeText >= fractionalText {
		t.Fatalf("formatted times are not chronologically sortable: %q >= %q", wholeText, fractionalText)
	}
	if len(wholeText) != len(fractionalText) {
		t.Fatalf("formatted times are not fixed width: %q and %q", wholeText, fractionalText)
	}
	for _, value := range []time.Time{whole, fractional} {
		parsed, err := parseTime(formatTime(value))
		if err != nil {
			t.Fatal(err)
		}
		if !parsed.Equal(value) {
			t.Fatalf("round trip = %s, want %s", parsed, value)
		}
	}
}

func TestParseTimeAcceptsLegacyRFC3339Nano(t *testing.T) {
	parsed, err := parseTime("2026-08-29T12:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, time.August, 29, 12, 0, 0, 0, time.UTC)
	if !parsed.Equal(want) {
		t.Fatalf("parsed = %s, want %s", parsed, want)
	}
}
