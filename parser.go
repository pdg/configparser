package main

import (
	"errors"
	"io"
	"strings"
)

type Directive struct {
	Name          string
	Arguments     []Argument
	Subdirectives []*Directive
}

type Argument string

func Parse(r io.Reader) ([]*Directive, error) {
	return ParseDirectives(NewTokenizer(r))
}

func ParseDirectives(z *Tokenizer) ([]*Directive, error) {

	var ds []*Directive

	for t, err := z.Next(); err != io.EOF; t, err = z.Next() {

		if t.Text == "{" || t.Text == "}" {
			return nil, errors.New("unexpected token: " + t.Text + "; expected start of directive")
		}

		if strings.TrimSpace(t.Text) == "" {
			continue
		}

		z.PutBack(t)

		d, err := ParseDirective(z)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}

	return ds, nil

}

func ParseDirective(z *Tokenizer) (*Directive, error) {

	d := &Directive{}

	err := z.SkipToNextToken()
	if err != nil {
		return nil, err
	}

	t, err := z.Next()
	if err != nil {
		return nil, err
	}

	d.Name = t.Text

	d.Arguments, err = ParseArguments(z)
	if err != nil && err != io.EOF {
		return nil, err
	}

	d.Subdirectives, err = ParseSubDirectives(z)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return d, nil

}

func ParseArguments(z *Tokenizer) ([]Argument, error) {

	var args []Argument

	for t, err := z.Next(); err != io.EOF; t, err = z.Next() {

		if t.Text == "{" || t.Text == "\n" {
			z.PutBack(t)
			break
		}

		args = append(args, Argument(t.Text))
	}

	return args, nil

}

func ParseSubDirectives(z *Tokenizer) ([]*Directive, error) {

	var ds []*Directive

	t, err := z.Next()

	if err != nil {
		return nil, err
	}

	// no sub directives
	if t.Text != "{" {
		z.PutBack(t)
		return ds, nil
	}

	for t, err := z.Next(); err != io.EOF; t, err = z.Next() {

		if err != nil {
			return nil, err
		}

		// make sure subdirective starts in new line
		if t.Text != "\n" {
			return nil, errors.New("unexpected token: " + t.Text + "; expected newline before new subdirective entry")
		}

		// skip linebreaks and whitespace characters
		err := z.Skip('\n', ' ', '\t')
		if err == io.EOF {
			return nil, errors.New("unexpected EOF in subdirectives block; closing curely brace missing")
		}
		if err != nil {
			return nil, err
		}

		// get next token
		t, err := z.Next()
		if err != nil {
			return nil, err
		}

		// break if subdirectives block ends
		if t.Text == "}" {
			break
		}

		// put back token if not end of subdirectives block
		z.PutBack(t)

		// parse subdirective
		d, err := ParseDirective(z)
		if err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil

}
