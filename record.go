package mmatcher

type Record struct {
	Id   string
	Atts []Atter
}

func (a *Record) IsMatch(b *Record, positions ...int) bool {
	if len(positions) <= 0 {
		positions = make([]int, len(a.Atts))
		for i := range positions {
			positions[i] = i
		}
	}
	e := make([]Atter, len(a.Atts))
	return a.IsMatchWithRanges(b, e, positions...)
}

func (a *Record) IsMatchWithRanges(b *Record, e []Atter, positions ...int) bool {
	if len(a.Atts) != len(b.Atts) || len(e) != len(a.Atts) {
		return false
	}
	if len(positions) <= 0 {
		positions = make([]int, len(a.Atts))
		for i := range positions {
			positions[i] = i
		}
	}
	matches := make([]bool, len(positions))
	for i, n := range positions {
		if n >= len(a.Atts) {
			return false
		}
		matches[i] = a.isMatchAt(b, e[n], n)
	}
	for _, m := range matches {
		if !m {
			return false
		}
	}
	return true
}

func (a *Record) isMatchAt(b *Record, e Atter, i int) bool {
	if i >= 0 && i < len(a.Atts) && i < len(b.Atts) {
		return a.Atts[i].Equal(b.Atts[i], e)
	}
	return false
}

type Records []Record

func (a *Record) MatchesAll(r Records, e ...Atter) []int {
	positions := make([]int, len(a.Atts))
	for i := range a.Atts {
		positions[i] = i
	}
	return a.Matches(r, positions, e...)
}

func (a *Record) Matches(r Records, positions []int, e ...Atter) (matches []int) {
	if len(e) <= 0 {
		e = make([]Atter, len(a.Atts))
	}
	for i, b := range r {
		if a.IsMatchWithRanges(&b, e, positions...) {
			matches = append(matches, i)
		}
	}
	return
}
