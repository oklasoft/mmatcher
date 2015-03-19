package mmatcher

type Record struct {
	Id   string
	Atts []Atter
}

func (a *Record) Matches(b *Record, positions ...int) bool {
	if len(positions) <= 0 {
		positions = make([]int, len(a.Atts))
		for i := range positions {
			positions[i] = i
		}
	}
	e := make([]Atter, len(a.Atts))
	return a.MatchesRanges(b, e, positions...)
}

func (a *Record) MatchesRanges(b *Record, e []Atter, positions ...int) bool {
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
		matches[i] = a.matchesAtt(b, e[n], n)
	}
	for _, m := range matches {
		if !m {
			return false
		}
	}
	return true
}

func (a *Record) matchesAtt(b *Record, e Atter, i int) bool {
	if i < len(a.Atts) && i < len(b.Atts) {
		return a.Atts[i].Equal(b.Atts[i], e)
	}
	return false
}
