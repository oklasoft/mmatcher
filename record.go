package mmatcher

type Record struct {
	Id   string
	Atts []Atter
}

func (a *Record) Matches(b *Record, atts ...int) bool {
	if len(atts) <= 0 {
		atts = make([]int, len(a.Atts))
		for i := range atts {
			atts[i] = i
		}
	}
	matches := make([]bool, len(atts))
	for i, n := range atts {
		if n >= len(a.Atts) || i >= len(b.Atts) {
			return false
		}
		matches[i] = a.matchesAtt(b, NumericAtt{}, n)
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
