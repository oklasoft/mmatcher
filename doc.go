// Copyright 2015 Stuart Glenn, OMRF. All rights reserved.
// Use of this code is governed by a 3 clause BSD style license
// Full license details in LICENSE file distributed with this software

// Package mmatcher package contains the structs & logic to find matching
// sets of Record based on some variable number of common attributes
//
// The basic idea is that one would get togehter a Records A, and Records B
// then be able for each in A find all the matches in B on some set of
// criteria (which attributes to check & optional +/- ranges therein)
package mmatcher
