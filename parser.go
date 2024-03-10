package configparser

import (
	"errors"
	"io"
	"strings"
)

type Directive struct {
	Name          string
	Arguments     []Argument
	Subdirectives Directives
}

type Argument string

type Directives []*Directive

func Parse(r io.Reader) (Directives, error) {
	dirs, err := parseDirectives(NewTokenizer(r))
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func (d Directives) Get(keys ...string) *Directive {

	if len(keys) == 0 {
		return nil
	}

	key, subkeys := keys[0], keys[1:]
	for _, d := range d {
		if d.Name == key {
			if len(subkeys) > 0 {
				return d.Subdirectives.Get(subkeys...)
			}
			return d
		}
	}

	return nil
}

func parseDirectives(z *Tokenizer) (Directives, error) {

	var ds []*Directive

	for t, err := z.Next(); err != io.EOF; t, err = z.Next() {

		if t.Text == "{" || t.Text == "}" {
			return Directives{}, errors.New("unexpected token: " + t.Text + "; expected start of directive")
		}

		if strings.TrimSpace(t.Text) == "" {
			continue
		}

		z.putBack(t)

		d, err := parseDirective(z)
		if err != nil {
			return Directives{}, err
		}
		ds = append(ds, d)
	}

	return Directives(ds), nil

}

func parseDirective(z *Tokenizer) (*Directive, error) {

	d := &Directive{}

	err := z.skipToNextToken()
	if err != nil {
		return nil, err
	}

	t, err := z.Next()
	if err != nil {
		return nil, err
	}

	d.Name = t.Text

	d.Arguments, err = parseArguments(z)
	if err != nil && err != io.EOF {
		return nil, err
	}

	d.Subdirectives, err = parseSubDirectives(z)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return d, nil

}

func parseArguments(z *Tokenizer) ([]Argument, error) {

	var args []Argument

	for t, err := z.Next(); err != io.EOF; t, err = z.Next() {

		if t.Text == "{" || t.Text == "\n" {
			z.putBack(t)
			break
		}

		args = append(args, Argument(t.Text))
	}

	return args, nil

}

func parseSubDirectives(z *Tokenizer) (Directives, error) {

	var ds Directives

	t, err := z.Next()

	if err != nil {
		return nil, err
	}

	// no sub directives
	if t.Text != "{" {
		z.putBack(t)
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
		err := z.skip('\n', ' ', '\t')
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
		z.putBack(t)

		// parse subdirective
		d, err := parseDirective(z)
		if err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil

}
