package gonv

import (
	"testing"
	"time"
)

func TestTimeParsing_DefaultFormat(t *testing.T) {
	s := time.Now().UTC().Format(DefaultTimeFormat)
	tm := Time(s)
	if tm.IsZero() {
		t.Fatalf("expected parsed time, got zero")
	}
}

func TestTimeUnix(t *testing.T) {
	// unix timestamp
	ts := int64(1_600_000_000)
	tm := Time(ts)
	if tm.Unix() != ts {
		t.Fatalf("expected unix %d, got %d", ts, tm.Unix())
	}
}
