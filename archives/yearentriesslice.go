// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package archives

import (
	"strconv"
)

type YearEntriesSlice []*YearEntries

// Len returns the length of the YearEntriesSlice.
func (ye YearEntriesSlice) Len() int {
	return len(ye)
}

// Less returns true if the value at i is newer than the value at j.
func (ye YearEntriesSlice) Less(i, j int) bool {
	first, err := strconv.Atoi(ye[j].Year)
	if err != nil {
		return true
	}

	second, err := strconv.Atoi(ye[i].Year)
	if err != nil {
		return true
	}

	return first < second
}

// Swap switches the elements at i and j.
func (ye YearEntriesSlice) Swap(i, j int) {
	ye[i], ye[j] = ye[j], ye[i]
}
