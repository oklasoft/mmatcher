// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/oklasoft/mmatcher/matcher"
)

func main() {
	case_file, err := os.Open(os.Args[1])
	if nil != err {
		log.Fatal(err)
	}
	defer case_file.Close()

	cases, err := matcher.NewRecordsFromCSV(case_file, 2)
	log.Print("Read cases:", cases)

	control_file, err := os.Open(os.Args[2])
	if nil != err {
		log.Fatal(err)
	}
	defer control_file.Close()

	controls, err := matcher.NewRecordsFromCSV(control_file, 2)
	log.Print("Read controls:", controls)

	all_matches := matcher.NewMatchSet()

	positions := []int{0, 1, 2}
	ranges := []matcher.Atter{matcher.TextAtt{}, matcher.TextAtt{}, matcher.NumericAtt{5}}

	for _, r := range cases {
		spots := r.Matches(controls, positions, ranges...)
		for _, i := range spots {
			all_matches.AddPair(matcher.NewPair(r.ID, controls[i].ID))
		}
	}

	log.Print("First pass lookikng for matches found:", all_matches)

	opti := all_matches.QuantityOptimized()

	log.Print("After optimizing we got:", opti)

	out := csv.NewWriter(os.Stdout)
	out.Write([]string{"case", "control"})

	for _, r := range cases {
		m := opti.MatchesFor(r.ID)
		if 0 == len(m) {
			continue
		}
		out.Write([]string{r.ID, m[0]})
	}
	out.Flush()

}
