package passgo

import (
	"strings"
	"sync"
	"testing"
)

func TestPasswordLength(t *testing.T) {
	g := New()
	for _, tc := range []struct {
		length, numbers, specials int
	}{
		{1, 0, 0},
		{2, 0, 0},
		{8, 0, 0},
		{12, 2, 1},
		{16, 4, 2},
		{5, 0, 0}, // odd word length
		{10, 10, 0},
		{10, 0, 10},
		{10, 5, 5},
	} {
		pw, err := g.Password(tc.length, tc.numbers, tc.specials)
		if err != nil {
			t.Errorf("Password(%d, %d, %d) unexpected error: %v", tc.length, tc.numbers, tc.specials, err)
			continue
		}
		if len(pw) != tc.length {
			t.Errorf("Password(%d, %d, %d) = %q, len %d; want %d", tc.length, tc.numbers, tc.specials, pw, len(pw), tc.length)
		}
	}
}

func TestPasswordComposition(t *testing.T) {
	g := New(WithCapitalization(0)) // disable caps for deterministic classification
	const numbers, specials, length = 3, 2, 12
	pw, err := g.Password(length, numbers, specials)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}

	var gotNums, gotSpecials, gotLetters int
	for _, c := range []byte(pw) {
		switch {
		case strings.IndexByte(string(DefaultNumbers), c) >= 0:
			gotNums++
		case strings.IndexByte(string(DefaultSpecialChars), c) >= 0:
			gotSpecials++
		default:
			gotLetters++
		}
	}
	if gotNums != numbers {
		t.Errorf("got %d digits, want %d (pw=%q)", gotNums, numbers, pw)
	}
	if gotSpecials != specials {
		t.Errorf("got %d specials, want %d (pw=%q)", gotSpecials, specials, pw)
	}
	if want := length - numbers - specials; gotLetters != want {
		t.Errorf("got %d letters, want %d (pw=%q)", gotLetters, want, pw)
	}
}

func TestPasswordErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		g                         *Generator
		length, numbers, specials int
	}{
		"zero length":       {New(), 0, 0, 0},
		"negative length":   {New(), -1, 0, 0},
		"negative numbers":  {New(), 8, -1, 0},
		"negative specials": {New(), 8, 0, -1},
		"over-allocated":    {New(), 4, 3, 2},
		"no consonants":     {New(WithConsonants(nil)), 8, 0, 0},
		"no vowels":         {New(WithVowels(nil)), 8, 0, 0},
		"numbers no set":    {New(WithNumbers(nil)), 8, 2, 0},
		"specials no set":   {New(WithSpecialChars(nil)), 8, 0, 2},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := tc.g.Password(tc.length, tc.numbers, tc.specials); err == nil {
				t.Errorf("expected an error, got nil")
			}
		})
	}
}

func TestNoCapitalization(t *testing.T) {
	g := New(WithCapitalization(0))
	for i := 0; i < 50; i++ {
		pw, err := g.Password(20, 0, 0)
		if err != nil {
			t.Fatalf("Password: %v", err)
		}
		if pw != strings.ToLower(pw) {
			t.Fatalf("expected no uppercase letters, got %q", pw)
		}
	}
}

func TestCapitalizationHappens(t *testing.T) {
	// With odds of 1, every letter should be upper-cased.
	g := New(WithCapitalization(1))
	pw, err := g.Password(20, 0, 0)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	if pw != strings.ToUpper(pw) {
		t.Errorf("with odds=1 every letter should be uppercase, got %q", pw)
	}
}

func TestCustomCharacterSets(t *testing.T) {
	g := New(
		WithConsonants([]byte("bcd")),
		WithVowels([]byte("ae")),
		WithCapitalization(0),
	)
	pw, err := g.Password(10, 0, 0)
	if err != nil {
		t.Fatalf("Password: %v", err)
	}
	const allowed = "bcdae"
	for _, c := range pw {
		if !strings.ContainsRune(allowed, c) {
			t.Errorf("password %q contains disallowed character %q", pw, c)
		}
	}
}

func TestConcurrentUse(t *testing.T) {
	g := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := g.Password(16, 3, 2); err != nil {
				t.Errorf("Password: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestUniqueness(t *testing.T) {
	// Cryptographically random passwords should essentially never collide.
	g := New()
	seen := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		pw, err := g.Password(12, 2, 1)
		if err != nil {
			t.Fatalf("Password: %v", err)
		}
		if _, dup := seen[pw]; dup {
			t.Fatalf("duplicate password generated: %q", pw)
		}
		seen[pw] = struct{}{}
	}
}

func BenchmarkPassword(b *testing.B) {
	g := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := g.Password(16, 3, 2); err != nil {
			b.Fatal(err)
		}
	}
}
