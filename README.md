passgo
======

[![CI](https://github.com/prinsmike/passgo/actions/workflows/ci.yml/badge.svg)](https://github.com/prinsmike/passgo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/prinsmike/passgo/v2.svg)](https://pkg.go.dev/github.com/prinsmike/passgo/v2)
[![Latest release](https://img.shields.io/github/v/release/prinsmike/passgo?sort=semver)](https://github.com/prinsmike/passgo/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/prinsmike/passgo)](go.mod)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

A human-readable ("pronounceable") password generator for Go. It builds
passwords from alternating consonants and vowels — easier to read and remember
than fully random strings — optionally mixed with digits and special
characters. It is based on [Pradeep Kishore Gowda's](https://code.activestate.com/recipes/410076-generate-a-human-readable-random-password-nicepass/) `nicepass.py`.

All randomness is drawn from `crypto/rand`, so the output is suitable for
security-sensitive use. A `Generator` holds no mutable state between calls and
is safe for concurrent use.

## Install

### Download a prebuilt binary

Prebuilt binaries for each [release](https://github.com/prinsmike/passgo/releases)
are attached as archives — no Go toolchain required. Grab the one for your
platform:

| OS      | Architecture | Asset                                  |
| ------- | ------------ | -------------------------------------- |
| Linux   | x86-64       | `passgo_<version>_linux_amd64.tar.gz`  |
| Linux   | ARM64        | `passgo_<version>_linux_arm64.tar.gz`  |
| macOS   | Intel        | `passgo_<version>_darwin_amd64.tar.gz` |
| macOS   | Apple silicon| `passgo_<version>_darwin_arm64.tar.gz` |
| Windows | x86-64       | `passgo_<version>_windows_amd64.zip`   |

For example, on Linux (x86-64):

```sh
VERSION=v2.0.0
curl -L -o passgo.tar.gz \
  "https://github.com/prinsmike/passgo/releases/download/${VERSION}/passgo_${VERSION}_linux_amd64.tar.gz"
tar -xzf passgo.tar.gz
sudo install passgo /usr/local/bin/
```

Each release also ships a `SHA256SUMS.txt`. To verify your download:

```sh
curl -LO "https://github.com/prinsmike/passgo/releases/download/${VERSION}/SHA256SUMS.txt"
sha256sum -c SHA256SUMS.txt --ignore-missing
```

### Install with the Go toolchain

Library:

```sh
go get github.com/prinsmike/passgo/v2
```

Command-line tool:

```sh
go install github.com/prinsmike/passgo/v2/cmd/passgo@latest
```

## Library usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/prinsmike/passgo/v2"
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
