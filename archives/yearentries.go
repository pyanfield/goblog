// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package archives

import (
	"github.com/pyanfield/goblog/blogs"
)

// This is used to simplify the Less function.
var mstrs map[string]int = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

// YearEntries is a list of entries and their associated year. It
// implements the sort interface for sorting the MonthEntries by date
// descending.
type YearEntries struct {
	// The name of the year.
	Year string

	// A list of MonthEntries for this year.
	Months []*MonthEntries
}

// Add appends the given BlogEntry to the Months list for the given
// month. It doesn't check to see if the given entry was actually in
// this year/month.
func (ye *YearEntries) Add(month string, e *blogs.BlogEntry) {
	// Make it if we don't have one.
	if ye.Months == nil {
		ye.Months = make([]*MonthEntries, 0, 0)
	}

	// Find the right month.
	var index int = -1
	for k, v := range ye.Months {
		if v.Month == month {
			index = k
			break
		}
	}

	if index == -1 {
		// We didn't find a month, so let's make it
		ye.Months = append(ye.Months, &MonthEntries{
			Month:   month,
			Entries: []*blogs.BlogEntry{e},
		})
	} else {
		ye.Months[index].Add(e)
	}
}

// Len returns the length of the YearEntries.
func (ye YearEntries) Len() int {
	return len(ye.Months)
}

// Less returns true if the value at i is newer than the value at j.
func (ye YearEntries) Less(i, j int) bool {
	return mstrs[ye.Months[j].Month] < mstrs[ye.Months[i].Month]
}

// Swap switches the elements at i and j.
func (ye YearEntries) Swap(i, j int) {
	ye.Months[i], ye.Months[j] = ye.Months[j], ye.Months[i]
}
