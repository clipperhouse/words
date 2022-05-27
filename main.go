package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

type config struct {
	In         *bufio.Reader
	HasIn      bool
	Delimiter  string
	Out        *bufio.Writer
	All        bool
	Lower      bool
	Upper      bool
	Diacritics bool
}

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
		Out: bufio.NewWriter(os.Stdout),
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435
	if piped {
		c.In = bufio.NewReader(os.Stdin)
		c.HasIn = true
	}

	flag.Parse()
	c.All = *all
	c.Lower = *lower
	c.Upper = *upper
	c.Diacritics = *diacritics

	if len(*delimiter) == 0 {
		c.Delimiter = "\n"
	} else {
		d, err := strconv.Unquote(`"` + *delimiter + `"`) // https://stackoverflow.com/a/59952849
		if err != nil {
			return nil, fmt.Errorf("couldn't parse delimiter %q: %v", *delimiter, err)

		}
		c.Delimiter = d
	}

	return c, nil
}

func writeWords(c *config) error {
	first := true
	sc := words.NewScanner(c.In)

	var transforms []transform.Transformer
	if c.Lower {
		transforms = append(transforms, transformer.Lower)
	}
	if c.Upper {
		transforms = append(transforms, transformer.Upper)
	}
	if c.Diacritics {
		transforms = append(transforms, transformer.Diacritics)
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
	log.Fatalln(err)
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
