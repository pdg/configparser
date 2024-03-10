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
		{"distribution", 0, Literal},
		{"debian", 14, Literal},
		{"stable", 21, Literal},
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
		# comment
	`

	wants := []*Token{
		{"\n", 0, Linebreak},
		{"distribution", 4, Literal},
		{"döbian", 18, Quoted},
		{"{", 26, Literal}, {"\n", 27, Linebreak}, {"\n", 31, Linebreak},
		{"suite", 35, Literal},
		{"stable", 42, Quoted}, {"\n", 49, Linebreak},
		{"architecture", 53, Literal},
		{"amd64 and more", 67, Quoted}, {"\n", 82, Linebreak}, {"\n", 83, Linebreak},
		{"repository", 87, Literal},
		{"{", 98, Literal}, {"\n", 99, Linebreak},
		{"security", 104, Literal}, {"\n", 112, Linebreak},
		{"backports", 117, Literal}, {"\n", 126, Linebreak},
		{"updates", 131, Literal}, {"\n", 138, Linebreak},
		{"}", 142, Literal}, {"\n", 143, Linebreak}, {"\n", 144, Linebreak},
		{"}", 147, Literal}, {"\n", 148, Linebreak},
		{" comment", 151, Comment}, {"\n", 160, Linebreak},
	}

	z := NewTokenizer(strings.NewReader(input))
	z.SkipComments = false

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
		{"a", 0, Literal},
		{"b", 2, Literal},
		{"\n", 3, Linebreak},
		{"c", 4, Literal},
		{"d", 6, Literal},
		{"{", 8, Literal},
		{"\n", 9, Linebreak},
		{"\n", 10, Linebreak},
		{"e", 13, Quoted},
		{"\n", 15, Linebreak},
		{"\n", 16, Linebreak},
		{"}", 17, Literal},
		{"\n", 18, Linebreak},
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

	err := z.skip('\n', ' ', '\t')
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	got, err := z.Next()
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	want := &Token{"a", 4, Literal}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}

}
