package main

import (
	"bufio"
	"io"
)

type Token struct {
	Text string
}

type Tokenizer struct {
	r         *bufio.Reader
	Seps      []rune
	prevToken []*Token
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:    bufio.NewReader(r),
		Seps: []rune{' ', '\t'}} // default seperator characters
}

func (z *Tokenizer) IsSep(r rune) bool {
	for _, sep := range z.Seps {
		if r == sep {
			return true
		}
	}
	return false
}

func (z *Tokenizer) Peako() (rune, error) {

	// read rune
	ru, i, err := z.r.ReadRune()
	if err != nil && err != io.EOF {
		return 0, err
	}

	// unread rune
	if i != 0 {
		err = z.r.UnreadRune()
		if err != nil {
			return 0, err
		}
	}

	return ru, nil

}

func (z *Tokenizer) readQutedString() (string, error) {

	var s string
	for ru, _, err := z.r.ReadRune(); err != io.EOF; ru, _, err = z.r.ReadRune() {

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

func (z *Tokenizer) SkipToContent() error {

	for ru, _, err := z.r.ReadRune(); err != io.EOF; ru, _, err = z.r.ReadRune() {

		if err != nil {
			return err
		}

		if !z.IsSep(ru) && ru != '\n' {
			err = z.r.UnreadRune()
			if err != nil {
				return err
			}
			break
		}

	}

	return nil
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

	for ru, _, err := z.r.ReadRune(); err != io.EOF; ru, _, err = z.r.ReadRune() {

		// read error
		if err != nil {
			return nil, err
		}

		// quoted string
		if ru == '"' && len(t.Text) == 0 {
			str, err := z.readQutedString()
			if err != nil {
				return nil, err
			}
			t.Text = str
			break
		}

		// seperator character
		if z.IsSep(ru) {
			if t.Text == "" {
				continue
			}
			break
		}

		// end of line
		if ru == '\n' {

			if t.Text == "" {
				t.Text = string(ru)
				break
			}

			err = z.r.UnreadRune()
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

	return t, nil

}
