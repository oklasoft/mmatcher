// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package matcher

import (
	"fmt"
	"log"
	"math"
)

//A Pair is the basic thing that makes up a match
type Pair struct {
	a string
	b string
}

//NewPair creates a new pairing between to items
func NewPair(a, b string) Pair {
	return Pair{a: a, b: b}
}

//Eql returns true if the two pairs are the same
//A Pair is the same regardless of the order of the elements
func (a Pair) Eql(b Pair) bool {
	return (a.a == b.a || a.a == b.b) && (a.b == b.a || a.b == b.b)
}

type matches []string

type match struct {
	m   matches
	isA bool
}

func newMatch(l int, is bool) *match {
	return &match{m: make(matches, l), isA: is}
}

func (m *match) copy() *match {
	m2 := newMatch(m.len(), m.isA)
	copy(m2.m, m.m)
	return m2
}

func (m *match) len() int {
	return len(m.m)
}

func (m *match) append(b string) {
	m.m = append(m.m, b)
}

func (m *match) delete(b string) bool {
	i := m.m.IndexOf(b)
	if i >= 0 {
		m.m = append(m.m[:i], m.m[i+1:]...)
		return true
	}
	return false
}

func (m match) String() string {
	return fmt.Sprintf("%v", m.m)
}

func (m matches) IndexOf(t string) int {
	for i, v := range m {
		if t == v {
			return i
		}
	}
	return -1
}

//MatchSet represents a collection of Pairs more or less
type MatchSet struct {
	pairs map[string]*match
}

func (m MatchSet) String() (s string) {
	for k, v := range m.pairs {
		if v.isA {
			s += fmt.Sprintf("%s(%t):%v ", k, v.isA, v.m)
		}
	}
	return s
}

//Copy returns a new copy of the MatchSet
func (m *MatchSet) Copy() (n MatchSet) {
	n = NewMatchSet()
	for k, v := range m.pairs {
		n.pairs[k] = v.copy()
	}
	return
}

//NewMatchSet creates a new MatchSet collection
func NewMatchSet() MatchSet {
	return MatchSet{make(map[string]*match)}
}

//AddPair adds a new pair of matched items to the collection
func (m *MatchSet) AddPair(p Pair) {
	if _, ok := m.pairs[p.a]; !ok {
		m.pairs[p.a] = newMatch(0, true)
	}
	(m.pairs[p.a]).append(p.b)
	(m.pairs[p.a]).isA = true
	if _, ok := m.pairs[p.b]; !ok {
		m.pairs[p.b] = newMatch(0, false)
	}
	(m.pairs[p.b]).append(p.a)
}

//RemovePair takes a pair of matched items out of the collection if its there
func (m *MatchSet) RemovePair(p Pair) {
	m.delete(p.a, p.b)
	m.delete(p.b, p.a)
}

//Purge takes a single ID, then removes any and all Pairs that contain at least half of it
func (m *MatchSet) Purge(i string) {
	if v, ok := m.pairs[i]; ok {
		for _, b := range v.m {
			//defer used, since RemovePair will alter the slice v & its range
			//this can lead to items being NOT purged, with defer we make sure
			//we loop through all the items. Defer used in place of copying v for "fun"
			defer m.RemovePair(NewPair(i, b))
		}
	}
}

func (m *MatchSet) delete(a, b string) {
	if v, ok := m.pairs[a]; ok {
		if v.delete(b) {
			if v.len() <= 0 {
				delete(m.pairs, a)
			}
		}
	}
}

//NumPairs returns the number of total pairs/matches in this collection
func (m *MatchSet) NumPairs() (l int) {
	for _, v := range m.pairs {
		l += v.len()
	}
	return l / 2
}

func (m *MatchSet) fewestPairs() (t string) {
	min := math.MaxInt32
	for k, v := range m.pairs {
		if v.len() < min {
			t = k
			min = v.len()
		}
	}
	return
}

func (m *MatchSet) mostPairsOf(t matches) (r string) {
	max := 0
	for _, p := range t {
		if m.pairs[p].len() > max {
			max = m.pairs[p].len()
			r = p
		}
	}
	return
}

func (m *MatchSet) minMaxPair() (p Pair) {
	p.a = m.fewestPairs()
	if "" != p.a {
		p.b = m.mostPairsOf(m.pairs[p.a].m)
	}
	return p
}

//QuantityOptimized returns an optimized matchset containing only a single
//pair per item. It attempts to get the largest number of possible pairs without
//duplicating any single item
func (m *MatchSet) QuantityOptimized(allowed ...int) (n MatchSet) {
	if len(allowed) <= 0 {
		allowed = []int{1}
	}
	n = NewMatchSet()
	if 0 == m.NumPairs() || allowed[0] <= 0 {
		return
	}
	c := m.Copy()
	nextSet := m.Copy()
	for p := c.minMaxPair(); "" != p.a && "" != p.b; p = c.minMaxPair() {
		if m.pairs[p.a].isA && m.pairs[p.b].isA {
			log.Fatal("Both cannot be A")
		} else if m.pairs[p.b].isA {
			tmp := p.a
			p.a = p.b
			p.b = tmp
		}
		c.Purge(p.a)
		c.Purge(p.b)
		nextSet.Purge(p.b)
		n.AddPair(p)
	}
	n.Add(nextSet.QuantityOptimized(allowed[0] - 1))
	return
}

func (m *MatchSet) Add(b MatchSet) {
	for k, v := range b.pairs {
		if v.isA {
			for _, p := range v.m {
				m.AddPair(NewPair(k, p))
			}
		}
	}
}

//MatchesFor returns a copy of the slice of matches for given identifeir.
func (m *MatchSet) MatchesFor(t string) (r matches) {
	if v, ok := m.pairs[t]; ok {
		r = make(matches, v.len())
		copy(r, v.m)
		return r
	}
	return r
}
