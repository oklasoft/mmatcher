// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package mmatcher

import "math"

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
	pairs map[string]matches
}

//Copy returns a new copy of the MatchSet
func (m *MatchSet) Copy() (n MatchSet) {
	n = NewMatchSet()
	for k, v := range m.pairs {
		n.pairs[k] = make([]string, len(v))
		copy(n.pairs[k], v)
	}
	return
}

//NewMatchSet creates a new MatchSet collection
func NewMatchSet() MatchSet {
	return MatchSet{make(map[string]matches)}
}

//AddPair adds a new pair of matched items to the collection
func (m *MatchSet) AddPair(p Pair) {
	m.pairs[p.a] = append(m.pairs[p.a], p.b)
	m.pairs[p.b] = append(m.pairs[p.b], p.a)
}

//RemovePair takes a pair of matched items out of the collection if its there
func (m *MatchSet) RemovePair(p Pair) {
	m.delete(p.a, p.b)
	m.delete(p.b, p.a)
}

//Purge takes a single ID, then removes any and all Pairs that contain at least half of it
func (m *MatchSet) Purge(i string) {
	if v, ok := m.pairs[i]; ok {
		for _, b := range v {
			m.RemovePair(NewPair(i, b))
		}
	}
}

func (m *MatchSet) delete(a, b string) {
	if v, ok := m.pairs[a]; ok {
		i := v.IndexOf(b)
		if i >= 0 {
			m.pairs[a] = append(v[:i], v[i+1:]...)
			if len(m.pairs[a]) <= 0 {
				delete(m.pairs, a)
			}
		}
	}
}

//NumPairs returns the number of total pairs/matches in this collection
func (m *MatchSet) NumPairs() (l int) {
	for _, v := range m.pairs {
		l += len(v)
	}
	return l / 2
}

func (m *MatchSet) fewestPairs() (t string) {
	min := math.MaxInt32
	for k, v := range m.pairs {
		if len(v) < min {
			t = k
			min = len(v)
		}
	}
	return
}

func (m *MatchSet) mostPairsOf(t []string) (r string) {
	max := 0
	for _, p := range t {
		if len(m.pairs[p]) > max {
			max = len(m.pairs[p])
			r = p
		}
	}
	return
}

//QuantityOptimized returns an optimized matchset containing only a single
//pair per item. It attempts to get the largest number of possible pairs without
//duplicating any single item
func (m *MatchSet) QuantityOptimized() (n MatchSet) {
	n = NewMatchSet()
	if 0 == m.NumPairs() {
		return
	}
	c := m.Copy()
	a := c.fewestPairs()
	b := c.mostPairsOf(m.pairs[a])
	c.Purge(a)
	c.Purge(b)
	n.AddPair(NewPair(a, b))
	return
}
