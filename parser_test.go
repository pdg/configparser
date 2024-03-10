package configparser

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestParseDirective(t *testing.T) {
	input := `distribution "debian" stable`
	wants := &Directive{Name: "distribution", Arguments: []Argument{"debian", "stable"}}

	tok := NewTokenizer(strings.NewReader(input))
	got, err := parseDirective(tok)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, wants) {
		t.Errorf("got %#v, wanted %#v", got, wants)
	}
}

func TestParseDirectives(t *testing.T) {

	input := `

    global # test

    distribution "debian" { 

      suite stable
      component "main" "contrib" "non free"

    }

  `
	wants := []*Directive{{
		Name:          "global",
		Arguments:     nil,
		Subdirectives: nil,
	}, {
		Name:      "distribution",
		Arguments: []Argument{"debian"},
		Subdirectives: []*Directive{{
			Name:          "suite",
			Arguments:     []Argument{"stable"},
			Subdirectives: nil,
		}, {
			Name:          "component",
			Arguments:     []Argument{"main", "contrib", "non free"},
			Subdirectives: nil,
		}},
	}}

	tok := NewTokenizer(strings.NewReader(input))
	got, err := parseDirectives(tok)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotj, _ := json.MarshalIndent(got, "", "  ")
	wantsj, _ := json.MarshalIndent(wants, "", "  ")

	if string(gotj) != string(wantsj) {
		t.Errorf("\ngot: %s\n\nwanted: %s", gotj, wantsj)
	}

}
