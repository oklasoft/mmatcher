package mmatcher

import "testing"

func TestMatch(t *testing.T) {
	a := &Record{Id: "a", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"green"}}}
	b := &Record{Id: "b", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"red"}}}
	if a.IsMatch(b) {
		t.Error("A should not match b on all attributes")
	}
	if !a.IsMatch(a) {
		t.Error("A should match itself on all attributes")
	}
	if !a.IsMatch(b, 0) {
		t.Error("A should match b on just first attribute")
	}
	if !b.IsMatch(a, []int{0, 1}...) {
		t.Error("B should match a on first two attributes")
	}
	if a.IsMatch(b, 2) {
		t.Error("A should not match B on third attribute")
	}
	if a.IsMatch(a, 100) {
		t.Error("A should not match itself with attribute out of bounds")
	}
	if !a.IsMatch(a, 0, 1, 2) {
		t.Error("A should match itself with all attributes specified")
	}
}

func TestMatchAtt(t *testing.T) {
	a := &Record{Id: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}, NumericAtt{8}, TextAtt{"yes"}}}
	b := &Record{Id: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
	if !a.isMatchAt(b, NumericAtt{4}, 1) {
		t.Error("A numeric att should match b with an epsilon")
	}
	if !a.isMatchAt(b, NumericAtt{3.2}, 1) {
		t.Error("A numeric att should match b with an epsilon")
	}
	if a.isMatchAt(b, NumericAtt{2}, 1) {
		t.Error("A numeric att should not match b with a small epsilon")
	}
	e := make([]Atter, len(a.Atts))
	e[1] = NumericAtt{2}
	if a.IsMatchWithRanges(b, e, 1) {
		t.Error("A should not match b with small range on one index")
	}
	if a.IsMatchWithRanges(b, e, 1, 2, 3) {
		t.Error("A should not match on multi indices b with small range on one index")
	}
	e[1] = NumericAtt{3.2}
	tests := map[int]bool{
		0: false,
		1: true,
		2: true,
		3: true,
	}
	for index, result := range tests {
		if result != a.isMatchAt(b, e[index], index) {
			t.Errorf("A:%v att:%v match? b:%v with a %v", a, index, b, e[index])
		}
	}
}

func TestMatchRange(t *testing.T) {
	a := &Record{Id: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}, NumericAtt{8}, TextAtt{"yes"}}}
	b := &Record{Id: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
	if a.IsMatch(b) {
		t.Error("A and b should not match")
	}
	e := make([]Atter, len(a.Atts))
	e[1] = NumericAtt{2}
	if a.IsMatchWithRanges(b, e, 1) {
		t.Error("A should not match b with small range on one index")
	}
	if a.IsMatchWithRanges(b, e, 1, 2, 3) {
		t.Error("A should not match on multi indices b with small range on one index")
	}
	e[1] = NumericAtt{3.2}
	i := []int{1}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with correct range:%v on %v", a, b, e, i)
	}
	i = []int{1, 2, 3}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with correct range:%v on %v", a, b, e, i)
	}
	e[1] = NumericAtt{3.0}
	if a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should not match b:%v with small range:%v on %v", a, b, e, i)
	}
	i = []int{2, 3}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with small range:%v on %v", a, b, e, i)
	}
}
