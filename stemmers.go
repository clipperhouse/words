package main

import (
	"github.com/clipperhouse/stemmer"
	"golang.org/x/text/transform"
)

var stemmerMap = map[string]transform.Transformer{
	"arabic":     stemmer.Arabic,
	"danish":     stemmer.Danish,
	"dutch":      stemmer.Dutch,
	"english":    stemmer.English,
	"finnish":    stemmer.Finnish,
	"french":     stemmer.French,
	"german":     stemmer.German,
	"hungarian":  stemmer.Hungarian,
	"irish":      stemmer.Irish,
	"italian":    stemmer.Italian,
	"norwegian":  stemmer.Norwegian,
	"porter":     stemmer.Porter,
	"portuguese": stemmer.Portuguese,
	"romanian":   stemmer.Romanian,
	"russian":    stemmer.Russian,
	"spanish":    stemmer.Spanish,
	"swedish":    stemmer.Swedish,
	"tamil":      stemmer.Tamil,
	"turkish":    stemmer.Turkish,
}
