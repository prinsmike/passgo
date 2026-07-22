package main

import (
	"strings"
	"testing"
)

const lengthFlag = "-length"

func TestRunGeneratesCount(t *testing.T) {
	var out strings.Builder
	if err := run([]string{"-count", "5", lengthFlag, "10"}, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("got %d passwords, want 5: %q", len(lines), out.String())
	}
	for _, l := range lines {
		if len(l) != 10 {
			t.Errorf("password %q has length %d, want 10", l, len(l))
		}
	}
}

func TestRunInvalidCount(t *testing.T) {
	var out strings.Builder
	if err := run([]string{"-count", "0"}, &out); err == nil {
		t.Error("expected error for count=0, got nil")
	}
}

func TestRunInvalidComposition(t *testing.T) {
	var out strings.Builder
	if err := run([]string{lengthFlag, "4", "-numbers", "3", "-specials", "3"}, &out); err == nil {
		t.Error("expected error when numbers+specials exceeds length, got nil")
	}
}

func TestRunNoCapitalization(t *testing.T) {
	var out strings.Builder
	if err := run([]string{"-odds", "0", lengthFlag, "20", "-numbers", "0", "-specials", "0"}, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	pw := strings.TrimSpace(out.String())
	if pw != strings.ToLower(pw) {
		t.Errorf("expected lowercase output with -odds 0, got %q", pw)
	}
}
