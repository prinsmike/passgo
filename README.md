passgo
======

[![CI](https://github.com/prinsmike/passgo/actions/workflows/ci.yml/badge.svg)](https://github.com/prinsmike/passgo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/prinsmike/passgo.svg)](https://pkg.go.dev/github.com/prinsmike/passgo)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

A human-readable ("pronounceable") password generator for Go. It builds
passwords from alternating consonants and vowels — easier to read and remember
than fully random strings — optionally mixed with digits and special
characters. It is based on [Pradeep Kishore Gowda's](https://code.activestate.com/recipes/410076-generate-a-human-readable-random-password-nicepass/) `nicepass.py`.

All randomness is drawn from `crypto/rand`, so the output is suitable for
security-sensitive use. A `Generator` holds no mutable state between calls and
is safe for concurrent use.

## Install

Library:

```sh
go get github.com/prinsmike/passgo
```

Command-line tool:

```sh
go install github.com/prinsmike/passgo/cmd/passgo@latest
```

## Library usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/prinsmike/passgo"
)

func main() {
	// Zero-config generator with sensible defaults.
	g := passgo.New()

	// 12 characters total: 2 digits, 1 special character, the rest letters.
	pass, err := g.Password(12, 2, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pass)
}
```

Override any of the defaults with options:

```go
g := passgo.New(
	passgo.WithConsonants([]byte("bcdfghjklmnpqrstvwxyz")),
	passgo.WithVowels([]byte("aeiou")),
	passgo.WithNumbers([]byte("0123456789")),
	passgo.WithSpecialChars([]byte("!@$#%&*-_.")),
	passgo.WithCapitalization(4), // 1-in-4 chance to capitalize each letter; 0 disables
)
```

## Command-line usage

```sh
$ passgo -count 3 -length 14 -numbers 3 -specials 2
```

Run `passgo -help` for the full list of flags:

```
-length int        total password length (default 12)
-numbers int       number of digits to include (default 2)
-specials int      number of special characters to include (default 1)
-count int         how many passwords to generate (default 1)
-odds int          1-in-N chance to capitalize each letter, 0 disables (default 4)
-consonants string consonant set
-vowels string     vowel set
-numchars string   digit set
-specialchars string special character set
```

## License

GPLv3. See [LICENSE](LICENSE).
