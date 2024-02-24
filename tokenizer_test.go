package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestTokenizerNext(t *testing.T) {

	input := `distribution  debian stable`
	wants := []*Token{
		{"distribution", 0},
		{"debian", 14},
		{"stable", 21},
	}

	z := NewTokenizer(strings.NewReader(input))
	for _, want := range wants {
		got, err := z.Next()
		if err != nil && err != io.EOF {
			t.Errorf("got	 %q, wanted nil or EOF", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, wanted %#v", got, want)
		}
	}
}

func TestTokenizerNextEOF(t *testing.T) {
	input := `directive`
	z := NewTokenizer(strings.NewReader(input))

	z.Next()               // distribution
	token, err := z.Next() // nil, EOF

	if token != nil {
		t.Errorf("got %q, wanted nil on eof", token)
	}

	if err != io.EOF {
		t.Errorf("got %q, wanted EOF", err)
	}

}

func TestTokenizerNextFull(t *testing.T) {

	input := `
   distribution "döbian" {
			
			suite "stable"
			architecture "amd64 and more"

			repository {
				security
				backports
				updates
			}

		}
	`

	wants := []*Token{
		{"\n", 0},
		{"distribution", 4},
		{"döbian", 18},
		{"{", 26}, {"\n", 27}, {"\n", 31},
		{"suite", 35},
		{"stable", 42}, {"\n", 49},
		{"architecture", 53},
		{"amd64 and more", 67}, {"\n", 82}, {"\n", 83},
		{"repository", 87},
		{"{", 98}, {"\n", 99},
		{"security", 104}, {"\n", 112},
		{"backports", 117}, {"\n", 126},
		{"updates", 131}, {"\n", 138},
		{"}", 142}, {"\n", 143}, {"\n", 144},
		{"}", 147}, {"\n", 148},
	}

	z := NewTokenizer(strings.NewReader(input))

	for _, want := range wants {
		got, err := z.Next()
		if err != nil && err != io.EOF {
			t.Errorf("got	 %q, wanted nil or EOF", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, wanted %#v", got, want)
		}
	}

}

func TestTokenizerNextShort(t *testing.T) {

	input := `a b
c d {

	"e"

}
`

	wants := []*Token{
		{"a", 0},
		{"b", 2},
		{"\n", 3},
		{"c", 4},
		{"d", 6},
		{"{", 8},
		{"\n", 9},
		{"\n", 10},
		{"e", 13},
		{"\n", 15},
		{"\n", 16},
		{"}", 17},
		{"\n", 18},
	}

	z := NewTokenizer(strings.NewReader(input))

	for _, want := range wants {
		got, err := z.Next()
		if err != nil && err != io.EOF {
			t.Errorf("got	 %q, wanted nil or EOF", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, wanted %#v", got, want)
		}
	}

}

func TestTokenizerSkip(t *testing.T) {

	input := `  
	a`

	z := NewTokenizer(strings.NewReader(input))

	err := z.Skip('\n', ' ', '\t')
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	got, err := z.Next()
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	want := &Token{"a", 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}

}
