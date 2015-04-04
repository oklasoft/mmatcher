// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/oklasoft/mmatcher/matcher"
	"gopkg.in/alecthomas/kingpin.v1"
)

func loadData(path string, skipHeader bool) matcher.Records {
	file, err := os.Open(path)
	if nil != err {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := matcher.NewRecordsFromCSV(file, skipHeader)
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

func version() string {
	return fmt.Sprintf("mmatcher - Multi Matcher 0.7.0 (20150404 %s)", build)
}

var (
	verbose      = kingpin.Flag("verbose", "Increase verbosity").Short('v').Bool()
	skipHeaders  = kingpin.Flag("skip-header", "Inputs have header line to be skipped, default is use everyline").Short('h').Bool()
	outFile      = kingpin.Flag("output", "Output file").Short('o').PlaceHolder("STDOUT").OpenFile(os.O_WRONLY|os.O_CREATE, 0660)
	outSep       = kingpin.Flag("out-separator", "Output field separator").Default(",").String()
	key          = kingpin.Arg("keys", "Keys to compare. A comma separated list of columns starting a 1, with optional :# +/- window").Required().String()
	case_file    = kingpin.Arg("case", "CSV file representing the cases").Required().ExistingFile()
	control_file = kingpin.Arg("controls", "CSV file representing the controls").Required().ExistingFile()
	build        string
)

func main() {
	kingpin.Version(version())
	kingpin.Parse()

	if nil == *outFile {
		outFile = &os.Stdout
	}

	positions, ranges := parseKeys(*key)
	cases := loadData(*case_file, *skipHeaders)
	controls := loadData(*control_file, *skipHeaders)

	all_matches := matcher.NewMatchSet()

	for _, r := range cases {
		spots := r.Matches(controls, positions, ranges...)
		for _, i := range spots {
			all_matches.AddPair(matcher.NewPair(r.ID, controls[i].ID))
		}
	}

	opti := all_matches.QuantityOptimized()

	out := csv.NewWriter(*outFile)
	sep, err := strconv.Unquote("'" + *outSep + "'")
	if nil != err {
		log.Fatal(err)
	}
	out.Comma = ([]rune(sep))[0]
	out.Write([]string{"case", "control"})

	for _, r := range cases {
		m := opti.MatchesFor(r.ID)
		if 0 == len(m) {
			out.Write([]string{r.ID})
			continue
		}
		line := []string{r.ID, m[0]}
		if *verbose {
			control := controls.Get(m[0])
			for _, p := range positions {
				line = append(line, r.Atts[p].String(), control.Atts[p].String())
			}
		}
		out.Write(line)
	}
	out.Flush()
	(*outFile).Close()

}
