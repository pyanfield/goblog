// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package archives

import (
	"github.com/pyanfield/goblog/blogs"
	"sort"
)

// DateEntries is a map that stores blog entries by year and month.
type DateEntries map[string]*YearEntries

// GetMost Recent returns up to the max most recent entries from the
// given (and hopefully sorted) list of YearEntries.
func GetMostRecent(y []*YearEntries, max int) []*blogs.BlogEntry {
	b := make([]*blogs.BlogEntry, 0, max)

	for _, year := range y {
		if len(b) == max {
			break
		}

		for _, month := range year.Months {
			if len(b) == max {
				break
			}

			for _, entry := range month.Entries {
				if len(b) == max {
					break
				}

				b = append(b, entry)
			}
		}
	}

	return b
}

// ParseBlogs creates a DateEntries from the given list of blogs.
func ParseBlogs(entries []*blogs.BlogEntry) DateEntries {
	t := make(DateEntries)

	for _, blog := range entries {
		year := blog.Created.Format("2006")
		month := blog.Created.Format("January")
		t.Add(year, month, blog)
	}

	return t
}

// Add stores the given BlogEntry under the given year and month.
func (de DateEntries) Add(year, month string, e *blogs.BlogEntry) {
	y, ok := de[year]
	if !ok {
		// We need to create it.
		de[year] = &YearEntries{
			Year: year,
			Months: []*MonthEntries{
				&MonthEntries{
					Month:   month,
					Entries: []*blogs.BlogEntry{e},
				},
			},
		}

		return
	}

	y.Add(month, e)
}

// Slice returns the Year, Month, and BlogEntries as a
// slice which is suitable for transformation in the
// templates. The year and months are sorted.
func (de DateEntries) Slice() YearEntriesSlice {
	s := make(YearEntriesSlice, 0, len(de))

	for _, y := range de {
		// Sort the months in the year.
		sort.Sort(y)

		// Sort each entry in each month.
		for _, m := range y.Months {
			sort.Sort(m)
		}

		// Add the year to our list.
		s = append(s, y)
	}

	// Now sort our year.
	sort.Sort(s)

	return s
}
