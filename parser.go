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

	z.SkipToContent()

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

		if t.Text == "" {
			break
		}

		// log.Println(args)
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

		// ENDLESS LOOP?

		if t.Text == "}" || t.Text == "" {
			break
		}

		d, err := ParseDirective(z)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}

	return ds, nil

}
