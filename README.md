# mmatcher - Multiple Matcher

[![Build Status](https://travis-ci.org/oklasoft/mmatcher.svg)](https://travis-ci.org/oklasoft/mmatcher)

## About

mmatcher is a simple CLI app to match cases & controls samples
for genetic studies

Matching is done between two CSV files containing an arbritrary number
of columns representing attributes on which to match. We make a poor
attempt to optimize the total number of matched case samples. A matched
control sample can only reported/used a single time.

## Installation

This is a Go app, so the standard go commands can be used to get & compile assuming
your go environment is setup

```shell
go get github.com/oklasoft/mmatcher
mmatcher --help
```

### Pre built binaries

If you trust me to build your binary for you the [releases](/oklasoft/mmatcher/releases) section has the
latest build for linux, mac & windows

## Usage

```shell
usage: mmatcher [<flags>] <keys> <case> <controls>

Flags:
  --help               Show help.
  -v, --verbose        Increase verbosity
  -h, --skip-header    Inputs have header line to be skipped, default is use everyline
  -o, --output=STDOUT  Output file
  -m, --matches=N      Allow up to N matches per case
  --out-separator=","  Output field separator
  --version            Show application version.

Args:
  <keys>      Keys to compare. A comma separated list of columns starting a 1, with optional :# +/- window
  <case>      CSV file representing the cases
  <controls>  CSV file representing the controls
```

### Input Files
mmatcher requires two CSV files to match each with one or more data columns upon which
to join. CSV files are required to have an identifer as the first column, each ID should be
unique within & between the two files. Data columns can then be anything you want after that.
Results will be best if the two files have their data columns in the same order. The first file
is the cases & the second the controls.

### Flags & Args

Matching keys/columns are specified via the *keys* arg by numbering the data columns, starting
at 1, after the ID column. The first column of data (the second overall column in the file) is
number 1, then 2, then 3, etc. All columns are matched if their contents exactly equal, unless
you specify a range for the column by appending :# to the key, where # is a number to use for
the +/- range. Of course that only really works if the data columns compared are numbers too.

Increased verbosity will cause output to include the data columns for mathches.
Normal output only includes the case ID & any matching control IDs. The data columns are listed
in order as specified by the *Keys* argument for the case, then each matching control.

The number of matches to be found per casae ID is 1 by default, you can try for more by using
the *-m* flag. We will try to find that many more additional matches per case. Control IDs are
still only used once & if matched against many cases, the case with the fewest matches at that
point gets the control.

The other flags will control output file or STDOUT, seperator (CSV or maybe tab) for output, etc.
By default input files are assumed to not have headers, so all lines are matched.

### Examples

```shell
mmatcher -m 10 -v 1,2,3:18,5 a.csv b.csv
```

Will find up to 10 matches for each line of a.csv from b.csv by comparing columns 1,2,5 & 3
as a number +/- 18 of each other. Output will be as a CSV to STDOUT of the matching IDs with
columns 1,2,3,5 of data for each ID in order

```shell
mmatcher -h -o matches.csv 1,2, a.csv b.csv
```

Will find & output just the IDs of one match from b to each from a saved as a CSV in the file matches.csv.
The first line of both a.csv & b.csv are skipped as headers.

```shell
mmatcher --out-separator "\t" -m 3 -v 2,3:18,4 a.csv b.csv
case	control	1	control	2	control	3
10474	33288	32185	32187	Female	Female	Female	Female	52	52	70	70	White	White	White	White
10481	26570	33749	32721	Female	Female	Female	Female	47	31	64	47	White	White	White	White
11390	33252	34430	34252	Female	Female	Female	Female	45	27	37	27	White	White	White	White
11427	34307	41638	34491	Female	Female	Female	Female	78	64	78	72	White	White	White	White
```

## License

Copyright 2015, Stuart Glenn, [Oklahoma Medical Research Foundation](https://omrf.org) (OMRF)

Distributed under a 3 clause BSD license. See [LICENSE](LICENSE) file
included in this software distrubtion for full details.
