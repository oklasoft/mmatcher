package mmatcher

import "testing"

func TestMatch(t *testing.T) {
	a := &Record{Id: "a", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"green"}}}
	b := &Record{Id: "b", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"red"}}}
	if a.Matches(b) {
		t.Error("A should not match b on all attributes")
	}
	if !a.Matches(a) {
		t.Error("A should match itself on all attributes")
	}
	if !a.Matches(b, 0) {
		t.Error("A should match b on just first attribute")
	}
	if !b.Matches(a, []int{0, 1}...) {
		t.Error("B should match a on first two attributes")
	}
	if a.Matches(b, 2) {
		t.Error("A should not match B on third attribute")
	}
	if a.Matches(a, 100) {
		t.Error("A should not match itself with attribute out of bounds")
	}
	if !a.Matches(a, 0, 1, 2) {
		t.Error("A should match itself with all attributes specified")
	}
}

func TestMatchRange(t *testing.T) {
	a := &Record{Id: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}, NumericAtt{8}, TextAtt{"yes"}}}
	b := &Record{Id: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
	if a.Matches(b) {
		t.Error("A and b should not match")
	}
	if !a.matchesAtt(b, NumericAtt{4}, 1) {
		t.Error("A numeric att should match b with an epsilon")
	}
	if a.matchesAtt(b, NumericAtt{2}, 1) {
		t.Error("A numeric att should not match b with a small epsilon")
	}
}
