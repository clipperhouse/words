package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	text := "Hello, goodbye"

	var w, werr bytes.Buffer
	r := strings.NewReader(text)
	c := &config{
		All: true,
		In:  r,
		Out: &w,
		Err: &werr,
	}

	err := write(c)
	if err != nil {
		t.Fatal(err)
	}

	expected := text
	got := w.String()

	if expected != got {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestUpperLower(t *testing.T) {
	text := "Hello, goodbye"

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Lower: true,
			In:    r,
			Out:   &w,
			Err:   &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "hellogoodbye"
		got := w.String()

		if expected != got {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Upper: true,
			In:    r,
			Out:   &w,
			Err:   &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "HELLOGOODBYE"
		got := w.String()

		if expected != got {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}
}

func TestDiacritics(t *testing.T) {
	text := "I am reading a résumé in Malmö"

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Diacritics: true,
			In:         r,
			Out:        &w,
			Err:        &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "IamreadingaresumeinMalmo"
		got := w.String()

		if expected != got {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Diacritics: true,
			Lower:      true,
			In:         r,
			Out:        &w,
			Err:        &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "iamreadingaresumeinmalmo"
		got := w.String()

		if expected != got {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}
}

func TestStemmer(t *testing.T) {
	text := "I am walking the Dogs"

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Stemmer: "english",
			In:      r,
			Out:     &w,
			Err:     &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "IamwalktheDog"
		got := w.String()

		if expected != got {
			t.Fatalf("expected %q, got %q", expected, got)
		}
	}

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Stemmer: "foo",
			In:      r,
			Out:     &w,
			Err:     &werr,
		}

		err := write(c)
		if err == nil {
			t.Fatalf("should have gotten an error for stem=%s", c.Stemmer)
		}
	}
}

func TestVersion(t *testing.T) {
	text := "I am walking the Dogs"

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Version: true,
			In:      r,
			Out:     &w,
			Err:     &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}
		got := werr.String()

		if !strings.HasPrefix(got, "Version: ") {
			t.Fatalf("expected 'Version: ', got %q", got)
		}
	}
}

func TestCount(t *testing.T) {
	text := "I am walking the Dogs"

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Count: true,
			In:    r,
			Out:   &w,
			Err:   &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "5"
		got := w.String()

		if expected != got {
			t.Fatalf("expected count of %q, got %q", expected, got)
		}
	}

	{
		var w, werr bytes.Buffer
		r := strings.NewReader(text)
		c := &config{
			Count: true,
			All:   true,
			In:    r,
			Out:   &w,
			Err:   &werr,
		}

		err := write(c)
		if err != nil {
			t.Fatal(err)
		}

		expected := "9"
		got := w.String()

		if expected != got {
			t.Fatalf("expected count of %q, got %q", expected, got)
		}
	}
}
