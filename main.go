package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/clipperhouse/uax29/iterators/filter"
	"github.com/clipperhouse/uax29/iterators/transformer"
	"github.com/clipperhouse/uax29/words"
	"golang.org/x/text/transform"
)

var all = flag.Bool("all", false, "include all tokens, such as whitespace and punctuation, not just 'words'")
var lower = flag.Bool("lower", false, "transform tokens to lower case")
var upper = flag.Bool("upper", false, "transform tokens to UPPER case")
var diacritics = flag.Bool("diacritics", false, "'flatten' / remove diacritic marks, such as accents, like açaí → acai")

var count = flag.Bool("count", false, "'count the number of words")

var delimiter = flag.String("delimiter", "", `separator to use between output tokens, default is "\n".
you can use escaped literals like "\t".`)
var stem = flag.String("stem", "", "language of a Snowball stemmer to apply to each token. options are:\narabic, danish, dutch, english, finnish, french, german, hungarian,\nirish, italian, norwegian, porter, portuguese, romanian, russian,\nspanish, swedish, tamil, turkish")

var v = flag.Bool("version", false, "print the current version and SHA")

type config struct {
	In  io.Reader
	Out writer
	Err writer

	Delimiter  string
	All        bool
	Lower      bool
	Upper      bool
	Diacritics bool
	Count      bool
	Stemmer    string
	Version    bool
}

var appName string = os.Args[0]
var version string
var commit string

func main() {
	c, err := getConfig()
	if err != nil {
		handle(err)
	}

	err = write(c)
	if err != nil {
		handle(err)
	}

	_, err = c.Out.WriteString("\n")
	if err != nil {
		handle(err)
	}
}

type writer interface {
	io.Writer
	io.StringWriter
}

func getConfig() (*config, error) {
	c := &config{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435
	if piped {
		c.In = os.Stdin
	}

	flag.Parse()

	c.All = *all
	c.Lower = *lower
	c.Upper = *upper
	c.Diacritics = *diacritics
	c.Stemmer = *stem

	if isFlagPassed("stem") {
		_, ok := stemmerMap[strings.ToLower(c.Stemmer)]
		if !ok {
			return nil, fmt.Errorf("unknown stemmer %q; type %q command for usage", c.Stemmer, appName)
		}
	}

	c.Count = *count
	c.Version = *v

	if isFlagPassed("delimiter") {
		c.Delimiter = *delimiter
	} else {
		c.Delimiter = `\n` // don't use "", we want to test the parsing in write()
	}

	return c, nil
}

func write(c *config) error {
	if c.Version {
		return printVersion(c)
	}

	if c.In == nil {
		return printUsage(c)
	}

	d, err := strconv.Unquote(`"` + c.Delimiter + `"`) // https://stackoverflow.com/a/59952849
	if err != nil {
		return fmt.Errorf("couldn't parse delimiter %q: %v", c.Delimiter, err)
	}
	c.Delimiter = d

	sc := words.NewScanner(c.In)

	var transformers []transform.Transformer
	if c.Diacritics {
		transformers = append(transformers, transformer.Diacritics)
	}
	if c.Lower {
		transformers = append(transformers, transformer.Lower)
	}
	if c.Upper {
		transformers = append(transformers, transformer.Upper)
	}
	if c.Stemmer != "" {
		stemmer, ok := stemmerMap[strings.ToLower(c.Stemmer)]
		if !ok {
			// a little redundant with above, but meh
			return fmt.Errorf("unknown stemmer %q; type %q command for usage", *stem, appName)
		}
		transformers = append(transformers, stemmer)
	}
	if len(transformers) > 0 {
		sc.Transform(transformers...)
	}

	if !c.All {
		sc.Filter(filter.Wordlike)
	}

	first := true
	count := 0
	for sc.Scan() {
		if c.Count {
			count++
			continue
		}

		if !first {
			_, err := c.Out.WriteString(c.Delimiter)
			if err != nil {
				return err
			}
		}
		first = false

		_, err := c.Out.Write(sc.Bytes())
		if err != nil {
			return err
		}
	}

	if sc.Err() != nil {
		return sc.Err()
	}

	if c.Count {
		_, err := c.Out.WriteString(strconv.Itoa(count))
		if err != nil {
			return err
		}
	}

	return nil
}

func handle(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func printUsage(c *config) error {
	flag.Usage()

	const message = `
Example:
  echo "Hello, 世界. Nice dog! 👍🐶" | words

Details:
  words accepts stdin, splits into one word (token) per line,
  and writes to stdout

  word boundaries are defined by Unicode, specifically UAX #29.
  by default, only tokens containing one or more letters,
  numbers, or symbols (as defined by Unicode) are returned;
  whitespace and punctuation tokens are omitted

`

	_, err := c.Err.WriteString(message)
	if err != nil {
		return err
	}

	err = printVersion(c)
	return err
}

func printVersion(c *config) error {
	v := fmt.Sprintf("Version: %s, SHA: %s\n", version, commit)
	_, err := c.Err.WriteString(v)
	return err
}

// https://stackoverflow.com/a/54747682
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
