package mmatcher

import (
	"fmt"
	"math"
)

type Atter interface {
	Equal(Atter, Atter) bool
}

type TextAtt struct {
	Val string
}

type NumericAtt struct {
	Val float64
}

func (a TextAtt) Equal(b Atter, e Atter) bool {
	v, ok := b.(TextAtt)
	return ok && a.Val == v.Val
}

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
