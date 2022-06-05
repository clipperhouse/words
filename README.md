`words` is a command which splits strings into individual words, as [defined by Unicode](https://unicode.org/reports/tr29/). It accepts text from stdin, and writes one word (token) per line to stdout.

Splitting words by the above standard is more likely to give consistent results; naïve splitting, such as by whitespace, is not precise. Further, the above standard is made for many languages and scripts.

`words` is intended to be helpful as part of a text pipeline with other transformation and search utilities.

### Install

Binaries are available on the [Releases page](https://github.com/clipperhouse/words/releases).

--

If you have a [Go installation](https://go.dev/doc/install), you can also `go install github.com/clipperhouse/words`.

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

### Motivation

Seems like this sort of primitive should exist!

If you’ve ever had to work with ‘words’ in an application, perhaps you made the naïve mistakes that I did. Splitting on whitespace should be good enough...oh but punctuation. Oh, also quotes. Oh, and other languages & scripts. It was always around 95% right, and 5% wrong is a big number.

The [Unicode standard](https://unicode.org/reports/tr29/) on which this tool is based handles the above well, across many types of text.

### Options

To see usage, just type the `words` command without arguments or input.

`-all`

By default, only ‘word’ tokens will be returned, omitting whitespace and punctuation tokens. Specify `-all` to return all tokens, i.e. 100% of the bytes.

`-delimiter`

A string separator to use between output tokens, default is `"\n"`. You can use escapes like `"\t"` for tab. It’s best to quote this parameter.

`-lower`

Lowercase the words.

`-upper`

UPPERCASE the words, if you must.

`-diacritics`

‘Flatten’ words with diacritics, such as accents. For example, açaí → acai.



