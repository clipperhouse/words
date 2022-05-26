package main

import (
	"bufio"
	"log"
	"os"

	"github.com/clipperhouse/uax29/iterators/filter"
	"github.com/clipperhouse/uax29/words"
)

func main() {
	fi, err := os.Stdin.Stat()
	handle(err)
	piped := (fi.Mode() & os.ModeCharDevice) == 0 // https://stackoverflow.com/a/43947435/70613

	out := bufio.NewWriter(os.Stdout)

	if !piped {
		out.WriteString("greetings\n")
		out.Flush()
		os.Exit(0)
	}

	delimiter := []byte("\n")

	in := bufio.NewReader(os.Stdin)

	first := true
	sc := words.NewScanner(in)
	sc.Filter(filter.Wordlike)

	for sc.Scan() {
		if !first {
			_, err := out.Write(delimiter)
			handle(err)
		}
		first = false
		_, err := out.Write(sc.Bytes())
		handle(err)
	}
	handle(sc.Err())

	_, err = out.WriteRune('\n')
	handle(err)
	handle(out.Flush())
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
