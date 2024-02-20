package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestTokenizerNext(t *testing.T) {

	input := `distribution "debian" stable`
	wants := []*Token{
		{"distribution"},
		{"debian"},
		{"stable"},
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
		distribution "debian" {

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
		{"\n"},
		{"distribution"},
		{"debian"},
		{"{"}, {"\n"}, {"\n"},
		{"suite"},
		{"stable"}, {"\n"},
		{"architecture"},
		{"amd64 and more"}, {"\n"}, {"\n"},
		{"repository"},
		{"{"}, {"\n"},
		{"security"}, {"\n"},
		{"backports"}, {"\n"},
		{"updates"}, {"\n"},
		{"}"}, {"\n"}, {"\n"},
		{"}"}, {"\n"},
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
