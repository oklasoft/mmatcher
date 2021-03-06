// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package matcher

import (
	"strings"
	"testing"
)

func TestMatch(t *testing.T) {
	a := &Record{ID: "a", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"green"}}}
	b := &Record{ID: "b", Atts: []Atter{TextAtt{"text"}, NumericAtt{1}, TextAtt{"red"}}}
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
	a := &Record{ID: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}, NumericAtt{8}, TextAtt{"yes"}}}
	b := &Record{ID: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
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

func TestMatchWrongSizes(t *testing.T) {
	a := &Record{ID: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}}}
	b := &Record{ID: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
	if a.IsMatchWithRanges(b, make([]Atter, len(a.Atts))) {
		t.Error("A longer record should not match a shorter record")
	}
	if a.IsMatchWithRanges(a, make([]Atter, 0)) {
		t.Error("A record cannot match itself is the []range is too short")
	}
	if !a.IsMatchWithRanges(a, make([]Atter, len(a.Atts))) {
		t.Error("A record should match itself completely")
	}
	if a.isMatchAt(a, TextAtt{}, len(a.Atts)+1) {
		t.Error("A record cannot equal even itself at a position past the record")
	}
	if a.isMatchAt(a, TextAtt{}, -1) {
		t.Error("A record cannot equal even itself at a position before the record")
	}
}

func TestMatchRange(t *testing.T) {
	a := &Record{ID: "a", Atts: []Atter{TextAtt{"first"}, NumericAtt{1}, NumericAtt{8}, TextAtt{"yes"}}}
	b := &Record{ID: "b", Atts: []Atter{TextAtt{"second"}, NumericAtt{4.2}, NumericAtt{8}, TextAtt{"yes"}}}
	if a.IsMatch(b) {
		t.Error("A and b should not match")
	}
	e := make([]Atter, len(a.Atts))
	e[0] = NumericAtt{2}
	if a.IsMatchWithRanges(b, e, 1) {
		t.Error("A should not match b with small range on one index")
	}
	if a.IsMatchWithRanges(b, e, 1, 2, 3) {
		t.Error("A should not match on multi indices b with small range on one index")
	}
	e[0] = NumericAtt{3.2}
	i := []int{1}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with correct range:%v on %v", a, b, e, i)
	}
	i = []int{1, 2, 3}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with correct range:%v on %v", a, b, e, i)
	}
	e[0] = nil
	e[1] = NumericAtt{3.0}
	if a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should not match b:%v with small range:%v on %v", a, b, e, i)
	}
	i = []int{2, 3}
	if !a.IsMatchWithRanges(b, e, i...) {
		t.Errorf("A:%v should match b:%v with small range:%v on %v", a, b, e, i)
	}
}

func TestMatchesAll(t *testing.T) {
	a := Records{
		Record{ID: "a1", Atts: []Atter{TextAtt{"red"}, NumericAtt{1}}},
		Record{ID: "a2", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{ID: "a3", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
	}
	b := Records{
		Record{ID: "b1", Atts: []Atter{TextAtt{"green"}, NumericAtt{5}}},
		Record{ID: "b2", Atts: []Atter{TextAtt{"red"}, NumericAtt{25}}},
		Record{ID: "b3", Atts: []Atter{TextAtt{"red"}, NumericAtt{35}}},
		Record{ID: "b4", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
		Record{ID: "b5", Atts: []Atter{TextAtt{"green"}, NumericAtt{31}}},
		Record{ID: "b6", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{ID: "b7", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
	}

	if m := a[0].MatchesAll(b); len(m) != 0 {
		t.Errorf("%v should not have found any matches in %v", a[0], b)
	}
	m := a[2].MatchesAll(b)
	if len(m) != 1 {
		t.Errorf("%v should have found one match in %v, but found %v", a[2], b, m)
	}
	if 3 != m[0] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[0])
	}

	m = a[1].MatchesAll(b)
	if len(m) != 2 {
		t.Errorf("%v should have found two matches in %v, but found %v", a[1], b, m)
	}
	if 5 != m[0] && 6 != m[1] {
		t.Errorf("%v should have found one at 5 & 6, but found it at %v", a[1], m)
	}

	e := []Atter{TextAtt{}, NumericAtt{5}}
	m = a[0].MatchesAll(b, e...)
	if 0 != len(m) {
		t.Errorf("%v should not have found any matches in %v using %v", a[0], b, e)
	}
	m = a[1].MatchesAll(b, e...)
	if 3 != len(m) {
		t.Errorf("%v should have found 3 in %v using %v, but found %v", a[0], b, e, m)
	}
	m = a[2].MatchesAll(b, e...)
	if 2 != len(m) {
		t.Errorf("%v should have found 2 in %v using %v, but found %v", a[0], b, e, m)
	}
}

func TestMatchesColumns(t *testing.T) {
	a := Records{
		Record{ID: "a0", Atts: []Atter{TextAtt{"red"}, NumericAtt{1}}},
		Record{ID: "a1", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{ID: "a2", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
	}
	b := Records{
		Record{ID: "b0", Atts: []Atter{TextAtt{"green"}, NumericAtt{5}}},
		Record{ID: "b1", Atts: []Atter{TextAtt{"red"}, NumericAtt{25}}},
		Record{ID: "b2", Atts: []Atter{TextAtt{"red"}, NumericAtt{35}}},
		Record{ID: "b3", Atts: []Atter{TextAtt{"green"}, NumericAtt{30}}},
		Record{ID: "b4", Atts: []Atter{TextAtt{"green"}, NumericAtt{31}}},
		Record{ID: "b5", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
		Record{ID: "b6", Atts: []Atter{TextAtt{"red"}, NumericAtt{20}}},
	}

	c := []int{10}
	if m := a[0].Matches(b, c); len(m) != 0 {
		t.Errorf("%v should have found 0 matches using column %v in %v, instead found %v", a[0], c, b, m)
	}

	c = []int{0}
	if m := a[0].Matches(b, c); len(m) != 4 {
		t.Errorf("%v should have found 4 matches in %v, instead found %v", a[0], b, m)
	}
	m := a[2].Matches(b, c)
	if len(m) != 3 {
		t.Errorf("%v should have found 3 matches in %v, but found %v", a[2], b, m)
	}
	if 0 != m[0] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[0])
	}
	if 3 != m[1] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[1])
	}
	if 4 != m[2] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[2])
	}
	c = []int{1}
	m = a[2].Matches(b, c)
	if len(m) != 1 {
		t.Errorf("%v should have found 1 matches in %v, but found %v", a[2], b, m)
	}
	if 3 != m[0] {
		t.Errorf("%v should have found one at 1, but found it at %v", a[2], m[0])
	}

	m = a[1].Matches(b, c)
	if len(m) != 2 {
		t.Errorf("%v should have found two matches in %v, but found %v", a[1], b, m)
	}
	if 5 != m[0] && 6 != m[1] {
		t.Errorf("%v should have found one at 5 & 6, but found it at %v", a[1], m)
	}

	e := []Atter{NumericAtt{5}}
	m = a[0].Matches(b, c, e...)
	if 1 != len(m) {
		t.Errorf("%v should have found 1 matches in %v using %v, but found %v", a[0], b, e)
	}
	m = a[1].Matches(b, c, e...)
	if 3 != len(m) {
		t.Errorf("%v should have found 3 in %v using %v, but found %v", a[0], b, e, m)
	}
	m = a[2].Matches(b, c, e...)
	if 4 != len(m) {
		t.Errorf("%v should have found 4 in %v using %v, but found %v", a[0], b, e, m)
	}
}

func TestCSVParsing(t *testing.T) {
	csv := `item,type,color,count
a1,m,red,25`
	r, err := NewRecordsFromCSV(strings.NewReader(csv), true)
	if err != nil {
		t.Error("Expected no error parsing, but got ", err)
	}
	if 1 != len(r) {
		t.Error("Expected 1 record from", r)
	}
	if 3 != len(r[0].Atts) {
		t.Error("Expedted 3 attributes from", r[0].Atts)
	}

	csv = `item,type,color,count
a1,f,red,15
a2,m,red,25`
	r, err = NewRecordsFromCSV(strings.NewReader(csv), true)
	if err != nil {
		t.Error("Expected no error parsing, but got ", err)
	}
	if 2 != len(r) {
		t.Error("Expected 2 record from", r)
	}
	if !(NumericAtt{15}).Equal(r[0].Atts[2], NumericAtt{}) {
		t.Error("Expected last attribute to be numeric equal to 15, but was not in", r[0].Atts[2])
	}

	csv = "item,type,color,count"
	r, err = NewRecordsFromCSV(strings.NewReader(csv), true)
	if err != nil {
		t.Error("Expected no error parsing, but got ", err)
	}
	if 0 != len(r) {
		t.Error("Expected 0 record from", r)
	}
	r, err = NewRecordsFromCSV(strings.NewReader(csv), false)
	if err != nil {
		t.Error("Expected no error parsing, but got ", err)
	}
	if 1 != len(r) {
		t.Error("Expected 1 record from", r)
	}
}

func TestCSVParsingWithCR(t *testing.T) {
	csv := "item,type,color,count\ra1,m,red,25\r"
	r, err := NewRecordsFromCSV(strings.NewReader(csv), true)
	if err != nil {
		t.Error("Expected no error parsing, but got ", err)
	}
	if 1 != len(r) {
		t.Fatal("Expected 1 record from", r, "in", csv)
	}
	if 3 != len(r[0].Atts) {
		t.Error("Expedted 3 attributes from", r[0].Atts)
	}

	csv = "item,type,color,count\ra1,f,red,15\ra2,m,red,25"
	r, err = NewRecordsFromCSV(strings.NewReader(csv), true)
	if err != nil {
		t.Fatal("Expected no error parsing, but got ", err)
	}
	if 2 != len(r) {
		t.Fatal("Expected 2 record from", r)
	}
	if !(NumericAtt{15}).Equal(r[0].Atts[2], NumericAtt{}) {
		t.Error("Expected last attribute to be numeric equal to 15, but was not in", r[0].Atts[2])
	}
}
