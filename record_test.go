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

func TestMatches(t *testing.T) {
	a := Records{
		Record{Id: "a1", Atts: []Atter{TextAtt{"red"}, NumericAtt{1}}},
		Record{Id: "a2", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{Id: "a3", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
	}
	b := Records{
		Record{Id: "b1", Atts: []Atter{TextAtt{"green"}, NumericAtt{5}}},
		Record{Id: "b2", Atts: []Atter{TextAtt{"red"}, NumericAtt{25}}},
		Record{Id: "b3", Atts: []Atter{TextAtt{"red"}, NumericAtt{35}}},
		Record{Id: "b4", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
		Record{Id: "b5", Atts: []Atter{TextAtt{"green"}, NumericAtt{31}}},
		Record{Id: "b6", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{Id: "b7", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
	}

	if m := a[0].Matches(b); len(m) != 0 {
		t.Errorf("%v should not have found any matches in %v", a[0], b)
	}
	m := a[2].Matches(b)
	if len(m) != 1 {
		t.Errorf("%v should have found one match in %v, but found %v", a[2], b, m)
	}
	if 3 != m[0] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[0])
	}

	m = a[1].Matches(b)
	if len(m) != 2 {
		t.Errorf("%v should have found two matches in %v, but found %v", a[1], b, m)
	}
	if 5 != m[0] && 6 != m[1] {
		t.Errorf("%v should have found one at 5 & 6, but found it at %v", a[1], m)
	}

	e := []Atter{TextAtt{}, NumericAtt{5}}
	m = a[0].Matches(b, e...)
	if 0 != len(m) {
		t.Errorf("%v should not have found any matches in %v using %v", a[0], b, e)
	}
	m = a[1].Matches(b, e...)
	if 3 != len(m) {
		t.Errorf("%v should have found 3 in %v using %v, but found %v", a[0], b, e, m)
	}
	m = a[2].Matches(b, e...)
	if 2 != len(m) {
		t.Errorf("%v should have found 2 in %v using %v, but found %v", a[0], b, e, m)
	}
}
