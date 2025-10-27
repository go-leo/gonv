package gonv

import (
	"testing"
	"time"
)

func TestDurationFromString(t *testing.T) {
	d := Duration("1h30m")
	if d != time.Hour+30*time.Minute {
		t.Fatalf("expected 1h30m, got %v", d)
	}
}

func TestDurationFromInt(t *testing.T) {
	// treat int as nanoseconds
	d := Duration(int64(2 * time.Second))
	if d != 2*time.Second {
		t.Fatalf("expected 2s, got %v", d)
	}
}

func TestDurationE_Error(t *testing.T) {
	_, err := DurationE("notaduration")
	if err == nil {
		t.Fatalf("expected error for invalid duration string")
	}
}
