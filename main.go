package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/clipperhouse/uax29/iterators/filter"
	"github.com/clipperhouse/uax29/iterators/transformer"
	"github.com/clipperhouse/uax29/words"
	"golang.org/x/text/transform"
)

var all = flag.Bool("all", false, "include all tokens, such as whitespace and punctuation, not just 'words'")
var lower = flag.Bool("lower", false, "transform tokens to lower case")
var upper = flag.Bool("upper", false, "transform tokens to UPPER case")
var diacritics = flag.Bool("diacritics", false, "'flatten' / remove diacritic marks, such as accents, like açaí → acai")

var delimiter = flag.String("delimiter", "", `separator to use between output tokens, default is "\n".
you can use escaped literals like "\t".`)
var stem = flag.String("stem", "", "language of a Snowball stemmer to apply to each token. options are:\narabic, danish, dutch, english, finnish, french, german, hungarian,\nirish, italian, norwegian, porter, portuguese romanian, russian,\nspanish, swedish, tamil, turkish")

type config struct {
	In         *bufio.Reader
	HasIn      bool
	Delimiter  string
	Out        *bufio.Writer
	All        bool
	Lower      bool
	Upper      bool
	Diacritics bool
	Stemmer    transform.Transformer
}

var appName string = os.Args[0]

func main() {
	config, err := getConfig()
	if err != nil {
		handle(err)
	}

	if !config.HasIn {
		printUsage()
		goto finish
	}

	err = writeWords(config)
	if err != nil {
		handle(err)
	}

finish:
	// Finish up
	_, err = config.Out.WriteString("\n")
	if err != nil {
		handle(err)
	}

	err = config.Out.Flush()
	if err != nil {
		handle(err)
	}
}

func getConfig() (*config, error) {
	c := &config{
		Out: bufio.NewWriterSize(os.Stdout, 64*1024),
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435
	if piped {
		c.In = bufio.NewReaderSize(os.Stdin, 64*1024)
		c.HasIn = true
	}

	flag.Parse()
	c.All = *all
	c.Lower = *lower
	c.Upper = *upper
	c.Diacritics = *diacritics

	if !isFlagPassed("delimiter") {
		c.Delimiter = "\n"
	} else {
		d, err := strconv.Unquote(`"` + *delimiter + `"`) // https://stackoverflow.com/a/59952849
		if err != nil {
			return nil, fmt.Errorf("couldn't parse delimiter %q: %v", *delimiter, err)

		}
		c.Delimiter = d
	}

	if isFlagPassed("stem") {
		stemmer, ok := stemmerMap[*stem]
		if !ok {
			return nil, fmt.Errorf("unknown stemmer %q; type %q command for usage", *stem, appName)
		}
		c.Stemmer = stemmer
	}

	return c, nil
}

func writeWords(c *config) error {
	first := true
	sc := words.NewScanner(c.In)

	var transforms []transform.Transformer
	if c.Diacritics {
		transforms = append(transforms, transformer.Diacritics)
	}
	if c.Lower {
		transforms = append(transforms, transformer.Lower)
	}
	if c.Upper {
		transforms = append(transforms, transformer.Upper)
	}
	if c.Stemmer != nil {
		transforms = append(transforms, c.Stemmer)
	}
	if len(transforms) > 0 {
		sc.Transform(transforms...)
	}

	if !c.All {
		sc.Filter(filter.Wordlike)
	}

	for sc.Scan() {
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

	return nil
}

func handle(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func printUsage() {
	message := "\nExample:\n  echo \"Hello, 世界. Nice dog! 👍🐶\" | words\n"
	message += "\nDetails:\n"
	message += "  words accepts stdin, splits into one word (token) per line,\n"
	message += "  and writes to stdout\n\n"
	message += "  word boundaries are defined by Unicode, specifically UAX #29.\n"
	message += "  by default, only tokens containing one or more letters,\n"
	message += "  numbers, or symbols (as defined by Unicode) are returned;\n"
	message += "  whitespace and punctuation tokens are omitted"

	flag.Usage()

	os.Stderr.WriteString(message)
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
