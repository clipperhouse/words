`words` is a command which splits strings into individual words, as [defined by Unicode](https://unicode.org/reports/tr29/).

It accepts text from stdin, and writes one word (token) per line to stdout.

### Install

```
go get github.com/clipperhouse/words
```

This requires a [Go installation](https://go.dev/doc/install). Note the bit about PATH in the instructions. (It's early days, we'll support installers like `brew` at some point.)

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

You can similarly use [`cat`](https://en.wikipedia.org/wiki/Cat_(Unix)) or [`curl`](https://curl.se/docs/manual.html) to stream data from a file or network, instead of `echo` above. You can also [pipe the output to a file](https://askubuntu.com/questions/420981/how-do-i-save-terminal-output-to-a-file).

### Options

To see options, just type `words` without arguments or input.

`-all`

By default, only 'word' tokens will be output, i.e., omitting whitespace or punctuation tokens. Specify `-all` to output all tokens, not just 'words'.

`-delimiter` string

a string separator to use between output tokens, default is `\n`. you can use escaped literals like `\t`. best to quote this parameter.
