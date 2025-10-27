package gonv

import "testing"

func TestBoolBasic(t *testing.T) {
	var b bool
	b = Bool[bool]("true")
	if !b {
		t.Fatalf("expected true, got false")
	}

	b = Bool[bool]("false")
	if b {
		t.Fatalf("expected false, got true")
	}

	b = Bool[bool](1)
	if !b {
		t.Fatalf("expected true for 1")
	}

	b = Bool[bool](0)
	if b {
		t.Fatalf("expected false for 0")
	}
}

func TestBoolE_Error(t *testing.T) {
	_, err := BoolE[bool]("notabool")
	if err == nil {
		t.Fatalf("expected error for invalid bool string")
	}
}
