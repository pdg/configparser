package main

import (
	"bufio"
	"io"
)

type TokenType int

const (
	Linebreak TokenType = iota
	Literal
	Quoted
	Comment
)

type Token struct {
	Text string
	Pos  int
	Type TokenType
}

type Tokenizer struct {
	r            *bufio.Reader
	Seps         []rune
	SkipComments bool
	prevToken    []*Token
	pos          int // current position in the tokenizer
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:            bufio.NewReader(r),
		Seps:         []rune{' ', '\t'}, // default seperator characters
		pos:          0,
		SkipComments: true,
	}
}

func (z *Tokenizer) isSep(ru rune) bool {
	for _, sep := range z.Seps {
		if ru == sep {
			return true
		}
	}
	return false
}

func (z *Tokenizer) readRune() (rune, int, error) {
	ru, i, err := z.r.ReadRune()
	if err == nil || err == io.EOF {
		z.pos++
	}
	return ru, i, err
}

func (z *Tokenizer) unreadRune() error {
	err := z.r.UnreadRune()
	if err == nil {
		z.pos = z.pos - 1
	}
	return err
}

func (z *Tokenizer) readQuoted() (string, error) {

	var s string
	for ru, _, err := z.readRune(); err != io.EOF; ru, _, err = z.readRune() {

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

func (z *Tokenizer) readComment() (string, error) {

	var s string
	for ru, _, err := z.readRune(); err != io.EOF; ru, _, err = z.readRune() {

		if err != nil {
			return "", err
		}

		if ru == '\n' {
			err = z.unreadRune()
			if err != nil {
				return "", err
			}
			break
		}

		s += string(ru)

	}

	return s, nil

}

func (z *Tokenizer) skipToNextToken() error {
	return z.skip(z.Seps...)
}

func (z *Tokenizer) skip(runes ...rune) error {

	for ru, _, err := z.readRune(); err != io.EOF; ru, _, err = z.readRune() {

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

		if err := z.unreadRune(); err != nil {
			return err
		} else {
			return nil
		}

	}

	return io.EOF

}

func (z *Tokenizer) putBack(t *Token) {
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

	for ru, _, err := z.readRune(); err != io.EOF; ru, _, err = z.readRune() {

		// read error
		if err != nil {
			return nil, err
		}

		// quoted string
		if ru == '"' && len(t.Text) == 0 {
			str, err := z.readQuoted()
			if err != nil {
				return nil, err
			}
			t.Text = str
			t.Type = Quoted
			break
		}

		// comment
		if ru == '#' && len(t.Text) == 0 {
			str, err := z.readComment()
			if err != nil && err != io.EOF {
				return nil, err
			}
			if z.SkipComments {
				continue
			}
			t.Text = str
			t.Type = Comment
			break
		}

		// seperator
		if z.isSep(ru) {
			if t.Text == "" {
				continue
			}
			break
		}

		// end of line
		if ru == '\n' {

			s++

			if t.Text == "" {
				t.Type = Linebreak
				t.Text = string(ru)
				break
			}

			err = z.unreadRune()
			if err != nil {
				return nil, err
			}

			break
		}

		// Literal
		t.Type = Literal
		t.Text += string(ru)

	}

	if t.Text == "" {
		return nil, io.EOF
	}

	t.Pos = z.pos - len([]rune(t.Text)) + s
	return t, nil

}
