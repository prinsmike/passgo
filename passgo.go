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

// Package passgo provides functions for generating human-readable passwords.

package passgo

import (
	"bytes"
	"errors"
	"math/rand"
	"time"
)

type Generator struct {
	Consonants     []byte
	Vowels         []byte
	Numbers        []byte
	SpecialChars   []byte
	Capitalize     bool
	CapitalizeOdds int
	buffer         bytes.Buffer
}

func (g *Generator) WriteChar(slice []byte) error {
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(len(slice))
	if g.Capitalize {
		g.buffer.WriteByte(g.ToUpper([]byte{slice[n]}))
	} else {
		g.buffer.WriteByte(slice[n])
	}
	return nil
}

func (g *Generator) ToUpper(char []byte) byte {
	var b byte
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(g.CapitalizeOdds)
	if n == g.CapitalizeOdds-1 {
		b = bytes.ToUpper(char)[0]
	} else {
		b = char[0]
	}
	return b
}

func (g *Generator) WriteWord(wlen int) error {
	for i := 0; i < wlen; i++ {
		if i%2 == 0 {
			g.WriteChar(g.Consonants)
		} else {
			g.WriteChar(g.Vowels)
		}
	}
	return nil
}

func (g *Generator) WriteNums(nLen int) error {
	for i := 0; i < nLen; i++ {
		g.WriteChar(g.Numbers)
	}
	return nil
}

func (g *Generator) WriteSpecialChars(sLen int) error {
	for i := 0; i < sLen; i++ {
		g.WriteChar(g.SpecialChars)
	}
	return nil
}

func (g *Generator) NewPassword(pLen, nLen, sLen int) (string, error) {
	if pLen <= 0 {
		err := errors.New("Passwords must be at least one character long.")
		return "", err
	}
	if len(g.Consonants) == 0 {
		err := errors.New("You must provide some consonants.")
		return "", err
	}
	if len(g.Vowels) == 0 {
		err := errors.New("You must provide some vowels.")
		return "", err
	}

	pLen = pLen - (nLen + sLen)

	if pLen%2 != 0 {
		g.WriteWord(pLen/2 + 1)
	} else {
		g.WriteWord(pLen / 2)
	}
	if len(g.Numbers) > 0 {
		g.WriteNums(nLen)
	}
	g.WriteWord(pLen / 2)
	if len(g.SpecialChars) > 0 {
		g.WriteSpecialChars(sLen)
	}

	return g.buffer.String(), nil
}

func NewGenerator(cons, vows, nums, specs []byte, caps bool, odds int) (g *Generator) {
	g = new(Generator)
	g.Consonants = cons
	g.Vowels = vows
	g.Numbers = nums
	g.SpecialChars = specs
	g.Capitalize = caps
	g.CapitalizeOdds = odds
	return
}
