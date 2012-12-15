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
	Vowels         []byte
	Consonants     []byte
	Numbers        []byte
	SpecialChars   []byte
	Capitalize     bool
	CapitalizeOdds int
}

func (g *Generator) GetChar(slice []byte) byte {
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(len(slice))

	return slice[n]
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

func (g *Generator) GetWord(wlen int) []byte {
	var wordslice []byte
	for i := 0; i < wlen; i++ {
		if i%2 == 0 {
			if g.Capitalize {
				wordslice = append(wordslice, g.ToUpper([]byte{g.GetChar(g.Vowels)}))
			} else {
				wordslice = append(wordslice, g.GetChar(g.Vowels))
			}
		} else {
			if g.Capitalize {
				wordslice = append(wordslice, g.ToUpper([]byte{g.GetChar(g.Consonants)}))
			} else {
				wordslice = append(wordslice, g.GetChar(g.Consonants))
			}
		}
	}
	return wordslice
}

func (g *Generator) GetNums(nlen int) []byte {
	var numslice []byte
	for i := 0; i < nlen; i++ {
		numslice = append(numslice, g.GetChar(g.Numbers))
	}
	return numslice
}

func (g *Generator) GetSpecialChars(clen int) []byte {
	var charslice []byte
	for i := 0; i < clen; i++ {
		charslice = append(charslice, g.GetChar(g.SpecialChars))
	}
	return charslice
}

func (g *Generator) GetPass(plen, nlen, clen int) ([]byte, error) {
	if plen <= 0 {
		err := errors.New("Passwords must be at least one character long.")
		return nil, err
	}
	if len(g.Consonants) == 0 {
		err := errors.New("You must provide some consonants.")
		return nil, err
	}
	if len(g.Vowels) == 0 {
		err := errors.New("You must provide some vowels.")
		return nil, err
	}
	var b bytes.Buffer
	if plen%2 != 0 {
		plen = plen + 1
	}
	b.Write(g.GetWord(plen / 2))
	if len(g.Numbers) > 0 {
		b.Write(g.GetNums(nlen))
	}
	b.Write(g.GetWord(plen / 2))
	if len(g.SpecialChars) > 0 {
		b.Write(g.GetSpecialChars(clen))
	}

	return b.Bytes(), nil
}
