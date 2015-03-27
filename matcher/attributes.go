// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package matcher

import (
	"fmt"
	"math"
)

// Atter is the interface to wrap comparing for equality between
// possible mixed string & numerics
type Atter interface {
	Equal(Atter, Atter) bool
	String() string
}

// A TextAtt is to store & compare string values for a Record
type TextAtt struct {
	Val string
}

// A NumericAtt is to store & compare number values for a Record
type NumericAtt struct {
	Val float64
}

// Equal returns true if strings a & b are in fact equal, e is ignored
func (a TextAtt) Equal(b Atter, e Atter) bool {
	v, ok := b.(TextAtt)
	return ok && a.Val == v.Val
}

// Equal returns true if numbers a & b are equal or within e (if e is NumericAtt)
// If e is not a NumericAtt, just a & b are compared for equality
func (a NumericAtt) Equal(b Atter, e Atter) bool {
	v, ok := b.(NumericAtt)
	if !ok {
		return false
	}
	epsilon, ok := e.(NumericAtt)
	if ok && epsilon.Val > 0 {
		return math.Abs(math.Abs(a.Val)-math.Abs(v.Val)) <= epsilon.Val
	}
	return a.Val == v.Val
}

func (a NumericAtt) String() string {
	return fmt.Sprintf("%v", a.Val)
}

func (a TextAtt) String() string {
	return a.Val
}
