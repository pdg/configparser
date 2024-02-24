package main

import (
	"bufio"
	"io"
)

type Token struct {
	Text string
	Pos  int
}

type Tokenizer struct {
	r         *bufio.Reader
	Seps      []rune
	prevToken []*Token
	pos       int // current position in the tokenizer.
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:    bufio.NewReader(r),
		Seps: []rune{' ', '\t'}, // default seperator characters
		pos:  0,
	}
}

func (z *Tokenizer) IsSep(ru rune) bool {
	for _, sep := range z.Seps {
		if ru == sep {
			return true
		}
	}
	return false
}

func (z *Tokenizer) ReadRune() (rune, int, error) {
	ru, i, err := z.r.ReadRune()
	if err == nil || err == io.EOF {
		z.pos++
	}
	return ru, i, err
}

func (z *Tokenizer) UnreadRune() error {
	err := z.r.UnreadRune()
	if err == nil {
		z.pos = z.pos - 1
	}
	return err
}

func (z *Tokenizer) readQuotedString() (string, error) {

	var s string
	for ru, _, err := z.ReadRune(); err != io.EOF; ru, _, err = z.ReadRune() {

		if err != nil {
			return "", err
		}

		// end of quoted string
		if ru == '"' {
			break
		}

		s += string(ru)

	}

	return s, nil

}

func (z *Tokenizer) SkipToNextToken() error {
	return z.Skip(z.Seps...)
}

func (z *Tokenizer) Skip(runes ...rune) error {

	for ru, _, err := z.ReadRune(); err != io.EOF; ru, _, err = z.ReadRune() {

		if err != nil {
			return err
		}

		var skip bool
		for _, r := range runes {
			if ru == r {
				skip = true
			}
		}

		if skip {
			continue
		}

		if err := z.UnreadRune(); err != nil {
			return err
		} else {
			return nil
		}

	}

	return io.EOF

}

func (z *Tokenizer) PutBack(t *Token) {
	z.prevToken = append(z.prevToken, t)
}

func (z *Tokenizer) Next() (*Token, error) {

	if l := len(z.prevToken); l > 0 {
		t := z.prevToken[l-1]
		z.prevToken = z.prevToken[:l-1]
		return t, nil
	}

	t := &Token{}

	s := -1

	for ru, _, err := z.ReadRune(); err != io.EOF; ru, _, err = z.ReadRune() {

		// read error
		if err != nil {
			return nil, err
		}

		// quoted string
		if ru == '"' && len(t.Text) == 0 {
			str, err := z.readQuotedString()
			if err != nil {
				return nil, err
			}
			t.Text = str
			break
		}

		// seperator
		if z.IsSep(ru) {
			if t.Text == "" {
				continue
			}
			break
		}

		// end of line
		if ru == '\n' {

			s++

			if t.Text == "" {
				t.Text = string(ru)
				break
			}

			err = z.UnreadRune()
			if err != nil {
				return nil, err
			}

			break
		}

		// normal characters
		t.Text += string(ru)

	}

	if t.Text == "" {
		return nil, io.EOF
	}

	t.Pos = z.pos - len([]rune(t.Text)) + s
	return t, nil

}
