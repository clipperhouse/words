package main

import (
	"bufio"
	"log"
	"os"

	"github.com/clipperhouse/uax29/iterators/filter"
	"github.com/clipperhouse/uax29/words"
)

type config struct {
	In        *bufio.Reader
	HasIn     bool
	Delimiter []byte
	Out       *bufio.Writer
}

func main() {
	config, err := getConfig()
	if err != nil {
		handle(err)
	}

	if !config.HasIn {
		message := "words handles piped input, splits into one word per line, and outputs to std out. use cat or echo to get started."
		config.Out.WriteString(message)
		goto finish
	}

	err = writeWords(config)
	if err != nil {
		handle(err)
	}

finish:
	// Finish up
	final := []byte("\n")
	_, err = config.Out.Write(final)
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
		Out:       bufio.NewWriter(os.Stdout),
		Delimiter: []byte("\n"),
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613
	if piped {
		c.In = bufio.NewReader(os.Stdin)
		c.HasIn = true
	}

	return c, nil
}

func writeWords(c *config) error {
	first := true
	sc := words.NewScanner(c.In)
	sc.Filter(filter.Wordlike)

	for sc.Scan() {
		if !first {
			_, err := c.Out.Write(c.Delimiter)
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
