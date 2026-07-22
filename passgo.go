/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// Package passgo generates human-readable ("pronounceable") passwords.
//
// Passwords are built from alternating consonants and vowels, which makes
// them easier to read and remember than fully random strings, optionally
// mixed with digits and special characters. All randomness is drawn from
// crypto/rand, so the output is suitable for security-sensitive use.
//
// A zero-configuration generator is available via New, which uses sane
// default character sets:
//
//	g := passgo.New()
//	pw, err := g.Password(12, 2, 1)
//
// A Generator is safe for concurrent use by multiple goroutines.
package passgo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Default character sets used when a Generator is created with New and no
// overriding options are supplied.
var (
	DefaultConsonants   = []byte("bcdfghjklmnpqrstvwxyz")
	DefaultVowels       = []byte("aeiou")
	DefaultNumbers      = []byte("0123456789")
	DefaultSpecialChars = []byte("!@$#%&*-_.")
)

// DefaultCapitalizeOdds is the default 1-in-N chance that any given letter is
// upper-cased when capitalization is enabled.
const DefaultCapitalizeOdds = 4

// Generator produces human-readable passwords from configurable character
// sets. Create one with New. The zero value is not usable; use New instead.
//
// A Generator holds no mutable state between calls, so a single Generator may
// be shared and used concurrently by multiple goroutines.
type Generator struct {
	consonants     []byte
	vowels         []byte
	numbers        []byte
	specialChars   []byte
	capitalize     bool
	capitalizeOdds int
}

// Option configures a Generator created by New.
type Option func(*Generator)

// WithConsonants sets the consonant set used to build word portions.
func WithConsonants(cons []byte) Option {
	return func(g *Generator) { g.consonants = cons }
}

// WithVowels sets the vowel set used to build word portions.
func WithVowels(vows []byte) Option {
	return func(g *Generator) { g.vowels = vows }
}

// WithNumbers sets the digit set used for the numeric portion.
func WithNumbers(nums []byte) Option {
	return func(g *Generator) { g.numbers = nums }
}

// WithSpecialChars sets the character set used for the special portion.
func WithSpecialChars(specs []byte) Option {
	return func(g *Generator) { g.specialChars = specs }
}

// WithCapitalization enables random capitalization. Each letter is upper-cased
// with a 1-in-odds probability. An odds value less than 1 disables
// capitalization.
func WithCapitalization(odds int) Option {
	return func(g *Generator) {
		if odds < 1 {
			g.capitalize = false
			return
		}
		g.capitalize = true
		g.capitalizeOdds = odds
	}
}

// New returns a Generator configured with the default character sets and
// capitalization enabled. Pass options to override any of the defaults.
func New(opts ...Option) *Generator {
	g := &Generator{
		consonants:     DefaultConsonants,
		vowels:         DefaultVowels,
		numbers:        DefaultNumbers,
		specialChars:   DefaultSpecialChars,
		capitalize:     true,
		capitalizeOdds: DefaultCapitalizeOdds,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Password generates a password of the given total length, containing numbers
// digits and specials special characters. The remaining length is filled with
// pronounceable word characters (alternating consonants and vowels).
//
// The layout is: word, digits, word, special characters.
//
// It returns an error if length is not positive, if numbers or specials is
// negative, if their sum exceeds length, or if the generator is missing a
// character set required to satisfy the request.
func (g *Generator) Password(length, numbers, specials int) (string, error) {
	switch {
	case length <= 0:
		return "", fmt.Errorf("passgo: length must be at least 1, got %d", length)
	case numbers < 0:
		return "", fmt.Errorf("passgo: numbers must not be negative, got %d", numbers)
	case specials < 0:
		return "", fmt.Errorf("passgo: specials must not be negative, got %d", specials)
	case numbers+specials > length:
		return "", fmt.Errorf("passgo: numbers+specials (%d) exceeds length (%d)", numbers+specials, length)
	case len(g.consonants) == 0:
		return "", fmt.Errorf("passgo: no consonants configured")
	case len(g.vowels) == 0:
		return "", fmt.Errorf("passgo: no vowels configured")
	case numbers > 0 && len(g.numbers) == 0:
		return "", fmt.Errorf("passgo: %d numbers requested but no number set configured", numbers)
	case specials > 0 && len(g.specialChars) == 0:
		return "", fmt.Errorf("passgo: %d specials requested but no special-character set configured", specials)
	}

	wordLen := length - numbers - specials
	// Split the word characters across two words, giving the extra character
	// to the first word when wordLen is odd.
	firstWord := wordLen/2 + wordLen%2
	secondWord := wordLen / 2

	var b strings.Builder
	b.Grow(length)

	if err := g.writeWord(&b, firstWord); err != nil {
		return "", err
	}
	if err := g.writeChars(&b, g.numbers, numbers, false); err != nil {
		return "", err
	}
	if err := g.writeWord(&b, secondWord); err != nil {
		return "", err
	}
	if err := g.writeChars(&b, g.specialChars, specials, false); err != nil {
		return "", err
	}

	return b.String(), nil
}

// writeWord writes n pronounceable characters, alternating consonants (on even
// positions) and vowels (on odd positions).
func (g *Generator) writeWord(b *strings.Builder, n int) error {
	for i := 0; i < n; i++ {
		set := g.consonants
		if i%2 == 1 {
			set = g.vowels
		}
		if err := g.writeChars(b, set, 1, true); err != nil {
			return err
		}
	}
	return nil
}

// writeChars writes n characters chosen uniformly at random from set. When
// letters is true and capitalization is enabled, each character may be
// upper-cased with a 1-in-capitalizeOdds probability.
func (g *Generator) writeChars(b *strings.Builder, set []byte, n int, letters bool) error {
	for i := 0; i < n; i++ {
		idx, err := randIndex(len(set))
		if err != nil {
			return err
		}
		c := set[idx]
		if letters && g.capitalize {
			up, err := g.maybeUpper(c)
			if err != nil {
				return err
			}
			c = up
		}
		b.WriteByte(c)
	}
	return nil
}

// maybeUpper returns the upper-case form of c with a 1-in-capitalizeOdds
// probability, and c unchanged otherwise.
func (g *Generator) maybeUpper(c byte) (byte, error) {
	roll, err := randIndex(g.capitalizeOdds)
	if err != nil {
		return c, err
	}
	if roll == 0 {
		return byte(strings.ToUpper(string(c))[0]), nil
	}
	return c, nil
}

// randIndex returns a uniformly distributed integer in [0, n) using a
// cryptographically secure source of randomness.
func randIndex(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("passgo: randIndex requires n > 0, got %d", n)
	}
	r, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, fmt.Errorf("passgo: reading random data: %w", err)
	}
	return int(r.Int64()), nil
}
