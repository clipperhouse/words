`words` is a command which splits strings into individual words, as [defined by Unicode](https://unicode.org/reports/tr29/). It accepts text from stdin, and writes one word (token) per line to stdout.

Splitting words by the above standard is more likely to give constent results; naive splitting, such as by whitespace, may not be precise. Further, the above standard is made for many languages and scripts.

`words` might be helpful as part of a text pipeline with other transformation and search utilities.

### Install

```
go get github.com/clipperhouse/words
```

This requires a [Go installation](https://go.dev/doc/install). Note the bit about PATH in the instructions. (It’s early days, we'll support installers like `brew` at some point.)

### Example

```
echo "Hello! Is this a puppy? 👍🐶" | words
```

Result

```
Hello
Is
this
a
puppy
👍
🐶
```

You can similarly use [`cat`](https://en.wikipedia.org/wiki/Cat_(Unix)) or [`curl`](https://curl.se/docs/manual.html) to stream data from a file or network, instead of `echo` above. You might pipe the output to [`sed`](https://www.gnu.org/software/sed/manual/sed.html), [`awk`](https://en.wikipedia.org/wiki/AWK) or [to a file](https://askubuntu.com/questions/420981/how-do-i-save-terminal-output-to-a-file).

You’ll note that by default, it only outputs ‘words’, defined as any token containing Unicode letters, numbers or symbols; whitespace and punctution tokens are omitted. See options below.

### Options

To see options, just type `words` without arguments or input.

`-all`

By default, only ‘word’ tokens will be output, omitting whitespace and punctuation tokens. Specify `-all` to output all tokens, not just ‘words’. I.e. delegate filtering to downstream systems.

`-delimiter`

A string separator to use between output tokens, default is `"\n"`. You can use escapes like `"\t"` for tab. It’s best to quote this parameter.
