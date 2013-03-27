// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package archives

import (
	"github.com/pyanfield/goblog/blogs"
)

// MonthEntries is a list of entries and their associated month. It
// implements the sort interface for sorting the entries by date
// descending.
type MonthEntries struct {
	// The name of the month.
	Month string

	// A list of blog entries for this month.
	Entries []*blogs.BlogEntry
}

// Add appends the given BlogEntry to the Entries list. It doesn't
// check to see if the given entry was actually in this month.
func (me *MonthEntries) Add(e *blogs.BlogEntry) {
	if me.Entries == nil {
		me.Entries = make([]*blogs.BlogEntry, 0, 0)
	}

	me.Entries = append(me.Entries, e)
}

// Len returns the length of the MonthEntries.
func (me *MonthEntries) Len() int {
	return len(me.Entries)
}

// Less returns true if the entry at j is newer than the entry at i.
func (me *MonthEntries) Less(i, j int) bool {
	return me.Entries[j].Created.Before(me.Entries[i].Created)
}

// Swap switches the elements at i and j.
func (me *MonthEntries) Swap(i, j int) {
	me.Entries[i], me.Entries[j] = me.Entries[j], me.Entries[i]
}
