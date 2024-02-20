package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestParseArguments(t *testing.T) {
}

func TestParseArgument(t *testing.T) {
}

func TestParseDirective(t *testing.T) {
	input := `  distribution "debian" stable `
	wants := &Directive{Name: "distribution", Arguments: []Argument{"debian", "stable"}}

	tok := NewTokenizer(strings.NewReader(input))
	got, err := ParseDirective(tok)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, wants) {
		t.Errorf("got %#v, wanted %#v", got, wants)
	}
}

func TestParseDirectives(t *testing.T) {

	input := `

		global

		distribution "debian" {
			suite stable
			component "main" "contrib" "non-free"
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
			Arguments:     []Argument{"main", "contrib", "non-free"},
			Subdirectives: nil,
		}},
	}}

	tok := NewTokenizer(strings.NewReader(input))
	got, err := ParseDirectives(tok)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotj, _ := json.MarshalIndent(got, "", "  ")
	wantsj, _ := json.MarshalIndent(wants, "", "  ")

	if !reflect.DeepEqual(gotj, wantsj) {
		t.Errorf("\ngot: %s\n\nwanted: %s", gotj, wantsj)
	}

}
