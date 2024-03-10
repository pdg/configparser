package configparser

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
		{"{", 26, Literal}, {"\n", 27, Linebreak}, {"\n", 34, Linebreak},
		{"suite", 41, Literal},
		{"stable", 48, Quoted}, {"\n", 55, Linebreak},
		{"architecture", 62, Literal},
		{"amd64 and more", 76, Quoted}, {"\n", 91, Linebreak}, {"\n", 92, Linebreak},
		{"repository", 99, Literal},
		{"{", 110, Literal}, {"\n", 111, Linebreak},
		{"security", 120, Literal}, {"\n", 128, Linebreak},
		{"backports", 137, Literal}, {"\n", 146, Linebreak},
		{"updates", 155, Literal}, {"\n", 162, Linebreak},
		{"}", 169, Literal}, {"\n", 170, Linebreak}, {"\n", 171, Linebreak},
		{"}", 176, Literal}, {"\n", 177, Linebreak},
		{" comment", 182, Comment}, {"\n", 191, Linebreak},
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
		{"e", 14, Quoted},
		{"\n", 16, Linebreak},
		{"\n", 17, Linebreak},
		{"}", 18, Literal},
		{"\n", 19, Linebreak},
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

	want := &Token{"a", 5, Literal}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}

}
