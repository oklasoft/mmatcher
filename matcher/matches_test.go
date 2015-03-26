// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package mmatcher_test

import (
	"testing"

	. "github.com/oklasoft/mmatcher"
)

func TestPair(t *testing.T) {
	a := NewPair("A1", "B1")
	b := NewPair("A2", "B2")
	c := NewPair("B1", "A1")
	if a.Eql(b) {
		t.Error(a, "not expected to equal", b)
	}
	if !a.Eql(c) {
		t.Error(a, "expected to equal", c)
	}
}

func TestAddRemove(t *testing.T) {
	m := NewMatchSet()
	if 0 != m.NumPairs() {
		t.Error("A new set should be empty, but got", m.NumPairs())
	}

	m.AddPair(NewPair("A1", "B1"))
	if 1 != m.NumPairs() {
		t.Error("Length expected to be 1, but got", m.NumPairs())
	}

	m.AddPair(NewPair("A1", "B2"))
	if 2 != m.NumPairs() {
		t.Error("Length expected to be 2, but got", m.NumPairs())
	}
	m.AddPair(NewPair("A1", "B3"))
	m.AddPair(NewPair("A2", "B2"))
	m.AddPair(NewPair("A2", "B3"))
	m.AddPair(NewPair("A3", "B1"))
	m.AddPair(NewPair("A3", "B2"))

	if 7 != m.NumPairs() {
		t.Error("Length expected to be 7, but got", m.NumPairs())
	}

	m.RemovePair(NewPair("A2", "B2"))
	if 6 != m.NumPairs() {
		t.Error("Afer Removing 1, length should be 6, but was", m.NumPairs(), "in", m)
	}
	m.RemovePair(NewPair("A1", "B1"))
	if 5 != m.NumPairs() {
		t.Error("Afer Removing 2, length should be 5, but was", m.NumPairs(), "in", m)
	}
	m.RemovePair(NewPair("A3", "B2"))
	if 4 != m.NumPairs() {
		t.Error("Afer Removing 3, length should be 4, but was", m.NumPairs(), "in", m)
	}
	m.RemovePair(NewPair("X3", "X2"))
	if 4 != m.NumPairs() {
		t.Error("After removing a non existent the size should still be 4", m.NumPairs(), "in", m)
	}
	m.RemovePair(NewPair("A1", "X2"))
	if 4 != m.NumPairs() {
		t.Error("After removing a non existent the size should still be 4", m.NumPairs(), "in", m)
	}
	m.RemovePair(NewPair("B3", "A2"))
	if 3 != m.NumPairs() {
		t.Error("After removing a pair in oppopsite order there should be 3, but there were", m.NumPairs(), "in", m)
	}
	m.Purge("A1")
	if 1 != m.NumPairs() {
		t.Error("After a purge of A3, we should have only 1 left, but was", m.NumPairs(), "in", m)
	}
}

func TestQuantityOptimize(t *testing.T) {
	m := NewMatchSet()
	if o := m.QuantityOptimized(); 0 != o.NumPairs() {
		t.Error("Expected an empty set of matches from an empty matchset", o)
	}
	m.AddPair(NewPair("A1", "B1"))
	m.AddPair(NewPair("A1", "B2"))
	if o := m.QuantityOptimized(); 1 != o.NumPairs() {
		t.Error("Expected only 1 optimized pair from 2 pairs with 3 samples, got", o)
	}
	m.AddPair(NewPair("A1", "B3"))
	m.AddPair(NewPair("A2", "B2"))
	if o := m.QuantityOptimized(); 2 != o.NumPairs() {
		t.Error("Expected 2 optimized pairs from 4 pairs with 4 samples, got", o)
	}
	m.AddPair(NewPair("A2", "B3"))
	if o := m.QuantityOptimized(); 2 != o.NumPairs() {
		t.Error("Expected 2 optimized pairs from 4 pairs with 5 samples, got", o)
	}
	m.AddPair(NewPair("A3", "B1"))
	if o := m.QuantityOptimized(); 3 != o.NumPairs() {
		t.Error("Expected 3 optimized pairs, got", o)
	}
	m.AddPair(NewPair("A3", "B2"))

	o := m.QuantityOptimized()
	if 3 != o.NumPairs() {
		t.Error("After optimizing we should have 3, but we had", o.NumPairs(), "in", o)
	}
	if 7 != m.NumPairs() {
		t.Error("After making the optimized set, the original should be the same size still", m)
	}
}
