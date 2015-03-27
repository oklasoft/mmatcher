// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/oklasoft/mmatcher/matcher"
)

func loadData(path string) matcher.Records {
	file, err := os.Open(path)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := matcher.NewRecordsFromCSV(file)
	if nil != err {
		log.Fatal(err)
	}
	return data
}

func parseKeys(s string) ([]int, []matcher.Atter) {
	parts := strings.Split(s, ",")
	positions := make([]int, len(parts))
	ranges := make([]matcher.Atter, len(parts))
	for i, v := range parts {
		k := strings.Split(v, ":")
		p, err := strconv.ParseInt(k[0], 10, 32)
		if nil != err {
			log.Fatal(err)
		}
		positions[i] = int(p) - 1
		if len(k) > 1 {
			r, err := strconv.ParseFloat(k[1], 32)
			if nil != err {
				log.Fatal(err)
			}
			ranges[i] = matcher.NumericAtt{r}
		}
	}
	return positions, ranges
}

func main() {
	positions, ranges := parseKeys(os.Args[1])
	log.Print(positions)
	log.Print(ranges)
	cases := loadData(os.Args[2])
	controls := loadData(os.Args[3])

	all_matches := matcher.NewMatchSet()

	for _, r := range cases {
		spots := r.Matches(controls, positions, ranges...)
		for _, i := range spots {
			all_matches.AddPair(matcher.NewPair(r.ID, controls[i].ID))
		}
	}

	opti := all_matches.QuantityOptimized()

	out := csv.NewWriter(os.Stdout)
	out.Write([]string{"case", "control"})

	for _, r := range cases {
		m := opti.MatchesFor(r.ID)
		if 0 == len(m) {
			out.Write([]string{r.ID})
			continue
		}
		control := controls.Get(m[0])
		line := []string{r.ID, m[0]}
		for _, p := range positions {
			line = append(line, r.Atts[p].String(), control.Atts[p].String())
		}
		out.Write(line)
	}
	out.Flush()

}
