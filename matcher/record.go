// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package matcher

import (
	"encoding/csv"
	"io"
	"strconv"
)

// A Record holds a data to be matched based on attributes in Atts
type Record struct {
	ID   string
	Atts []Atter
}

// IsMatch returns true if Record a matches b exactly in columns given by positions
func (a *Record) IsMatch(b *Record, positions ...int) bool {
	if len(positions) <= 0 {
		positions = make([]int, len(a.Atts))
		for i := range positions {
			positions[i] = i
		}
	}
	e := make([]Atter, len(positions))
	return a.IsMatchWithRanges(b, e, positions...)
}

// IsMatchWithRanges returns true if Record a matches b in columns specified in
// positions. e is a slice of Atters to use for +/- range comparisons in columns
// of the same index
func (a *Record) IsMatchWithRanges(b *Record, e []Atter, positions ...int) bool {
	if len(a.Atts) != len(b.Atts) {
		return false
	}
	if len(positions) <= 0 {
		positions = make([]int, len(a.Atts))
		for i := range positions {
			positions[i] = i
		}
	}
	if len(positions) > len(e) {
		return false
	}
	matches := make([]bool, len(positions))
	for i, n := range positions {
		matches[i] = a.isMatchAt(b, e[i], n)
	}
	for _, m := range matches {
		if !m {
			return false
		}
	}
	return true
}

// isMatchAt returns true if single attribute column in i matches between
// a & b with given +/- range e
func (a *Record) isMatchAt(b *Record, e Atter, i int) bool {
	if i >= 0 && i < len(a.Atts) && i < len(b.Atts) {
		return a.Atts[i].Equal(b.Atts[i], e)
	}
	return false
}

// Records is just a slice of Record types
type Records []Record

//NewRecordsFromCSV parses an CSV formatted io.Reader to create
//Records for matching. We assume the first line is a header row which
//is skipped.
//TODO we should make this more robust with checking number of columns etc
func NewRecordsFromCSV(in io.Reader, skipHeader bool) (r Records, err error) {
	csv := csv.NewReader(in)
	lineno := 0

	for {
		lineno++
		line, err := csv.Read()
		if io.EOF == err {
			err = nil
			break
		} else if nil != err {
			return nil, err
		}
		if skipHeader && 1 == lineno {
			continue //skip header
		}
		a := []Atter{}
		for _, v := range line[1:] {
			n, err := strconv.ParseFloat(v, 64)
			if nil == err {
				a = append(a, NumericAtt{n})
			} else {
				a = append(a, TextAtt{v})
			}
		}
		r = append(r, Record{ID: line[0], Atts: a})
	}

	return r, nil
}

func (r *Records) Get(t string) Record {
	for _, v := range *r {
		if v.ID == t {
			return v
		}
	}
	return Record{}
}

// MatchesAll returns a slice containing the indices of r that match to a with
// the given +/- ranges in e
func (a *Record) MatchesAll(r Records, e ...Atter) []int {
	positions := make([]int, len(a.Atts))
	for i := range a.Atts {
		positions[i] = i
	}
	return a.Matches(r, positions, e...)
}

// Matches retruns a slice containing the indices from r that match to a at
// attributes in positions with any given +/- ranges in e
func (a *Record) Matches(r Records, positions []int, e ...Atter) (matches []int) {
	if len(e) <= 0 {
		e = make([]Atter, len(positions))
	}
	for i, b := range r {
		if a.IsMatchWithRanges(&b, e, positions...) {
			matches = append(matches, i)
		}
	}
	return
}
