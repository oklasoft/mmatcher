package mmatcher_test

import "testing"

import . "github.com/oklasoft/mmatcher"

func TestTextAttsEqual(t *testing.T) {
	ta := TextAtt{"hi"}
	ta2 := TextAtt{"hi"}
	if !ta.Equal(ta2, TextAtt{}) {
		t.Errorf("%s expected equal to %s", ta, ta2)
	}
	if ta.Equal(ta2, TextAtt{}) != ta2.Equal(ta, TextAtt{}) {
		t.Error("Expected to be the same")
	}
	ta2.Val = "nope"
	if ta2.Equal(ta, TextAtt{}) {
		t.Error("%s expected NOT to equal %s", ta2, ta)
	}
	if !ta.Equal(ta, TextAtt{}) {
		t.Error("TextAtt expected to euqal itself")
	}
}

func TestNumericAttsEqual(t *testing.T) {
	n1 := NumericAtt{0}
	if !n1.Equal(n1, NumericAtt{}) {
		t.Error("NumericAtt expected to equal itself with no epsilon")
	}
	if !n1.Equal(n1, NumericAtt{1}) {
		t.Error("NumericAtt expected to equal itself with epsilon")
	}
	if n1.Equal(TextAtt{"fail city"}, NumericAtt{}) {
		t.Error("NumericAtt should not equal non NumericAtt{}")
	}
	n2 := NumericAtt{10.0}
	if n1.Equal(n2, NumericAtt{}) {
		t.Errorf("%s should not equal %s without epsilon", n1, n2)
	}
	e := NumericAtt{10}
	if !n2.Equal(n1, e) {
		t.Errorf("%s should equal %s with epsilon %s", n2, n1, e)
	}
	e.Val = 9
	if n1.Equal(n2, e) {
		t.Errorf("%s should not equal %s with epsilon %s", n1, n2, e)
	}
	n1.Val = 15.3
	e.Val = 5.4
	if !n1.Equal(n2, e) {
		t.Errorf("%s should equal %s with epsilon %s", n1, n2, e)
	}
	e.Val = 5.29
	if n1.Equal(n2, e) {
		t.Error("%s should not equal %s with epsilon %s", n1, n2, e)
	}
}
