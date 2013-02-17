// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package tags

import (
	"github.com/icub3d/goblog/blogs"
	"sort"
)

// TagEntries is a map of TagEntry structures with some methods for
// easily adding blog entries. It also has the ability to export the
// entries as a list for further processing.
type TagEntries map[string]*TagEntry

// ParseBlogs builds a TagEntries from the given list of blogs.
func ParseBlogs(entries []*blogs.BlogEntry) TagEntries {
	t := make(TagEntries)

	for _, blog := range entries {
		t.Add(blog)
	}

	return t
}

// Slice returns the TagEntry structures with this TagEntries as a
// slice. The list is in sorted order.
func (te TagEntries) Slice() TagEntriesSlice {
	s := make(TagEntriesSlice, 0, len(te))

	for _, t := range te {
		s = append(s, t)
	}

	// Sort the slice.
	sort.Sort(s)

	return s
}

// Add links the given BlogEntry to all of it's tags.
func (te TagEntries) Add(e *blogs.BlogEntry) {
	for _, tag := range e.Tags {
		f, ok := te[tag]
		if !ok {
			// It wasn't found, so create one.
			te[tag] = &TagEntry{
				Name:    tag,
				Entries: []*blogs.BlogEntry{e},
			}
		} else {
			// Add it to the one we found.
			f.Add(e)
		}
	}
}
