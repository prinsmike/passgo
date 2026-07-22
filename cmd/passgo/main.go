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

// Command passgo generates human-readable passwords on the command line.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/prinsmike/passgo"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "passgo:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("passgo", flag.ContinueOnError)
	fs.SetOutput(out)

	length := fs.Int("length", 12, "total password length")
	numbers := fs.Int("numbers", 2, "number of digits to include")
	specials := fs.Int("specials", 1, "number of special characters to include")
	count := fs.Int("count", 1, "how many passwords to generate")
	odds := fs.Int("odds", passgo.DefaultCapitalizeOdds, "1-in-N chance to capitalize each letter (0 disables)")
	consonants := fs.String("consonants", string(passgo.DefaultConsonants), "consonant set")
	vowels := fs.String("vowels", string(passgo.DefaultVowels), "vowel set")
	numchars := fs.String("numchars", string(passgo.DefaultNumbers), "digit set")
	specialchars := fs.String("specialchars", string(passgo.DefaultSpecialChars), "special character set")

	fs.Usage = func() {
		fmt.Fprintln(out, "Usage: passgo [flags]")
		fmt.Fprintln(out, "\nGenerate human-readable (pronounceable) passwords.")
		fmt.Fprintln(out, "\nFlags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *count < 1 {
		return fmt.Errorf("count must be at least 1, got %d", *count)
	}

	g := passgo.New(
		passgo.WithConsonants([]byte(*consonants)),
		passgo.WithVowels([]byte(*vowels)),
		passgo.WithNumbers([]byte(*numchars)),
		passgo.WithSpecialChars([]byte(*specialchars)),
		passgo.WithCapitalization(*odds),
	)

	for i := 0; i < *count; i++ {
		pw, err := g.Password(*length, *numbers, *specials)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, pw)
	}
	return nil
}
