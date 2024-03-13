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

func TestFirst(t *testing.T) {
	input := `
	distribution "debian" { 
		suite stable {
			component "main" "contrib" "non free"
		}		
	}
	`

	config, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	got := config.First("distribution", "suite", "component")
	wants := &Directive{Name: "component", Arguments: []Argument{"main", "contrib", "non free"}, Subdirectives: nil}

	gotj, _ := json.MarshalIndent(got, "", "  ")
	wantsj, _ := json.MarshalIndent(wants, "", "  ")

	if string(gotj) != string(wantsj) {
		t.Errorf("\ngot: %s\n\nwanted: %s", gotj, wantsj)
	}

}

func TestAll(t *testing.T) {
	input := `
	node { 
		role "one"
		role "two"
	}
	node {
		role "three"
	}
	`

	config, err := Parse(strings.NewReader(input))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	got := config.All("node", "role")
	wants := Directives{
		&Directive{Name: "role", Arguments: []Argument{"one"}, Subdirectives: nil},
		&Directive{Name: "role", Arguments: []Argument{"two"}, Subdirectives: nil},
		&Directive{Name: "role", Arguments: []Argument{"three"}, Subdirectives: nil},
	}

	gotj, _ := json.MarshalIndent(got, "", "  ")
	wantsj, _ := json.MarshalIndent(wants, "", "  ")

	if string(gotj) != string(wantsj) {
		t.Errorf("\ngot: %s\n\nwanted: %s", gotj, wantsj)
	}

}
