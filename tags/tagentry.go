// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package tags

import (
	"github.com/icub3d/goblog/blogs"
)

// TagEntry is a Tag and it's associated entries. It is used as a
// storage mechanism for the tags page.
type TagEntry struct {
	// The name of the Tag.
	Name string

	// The list of Entries associated with this tag.
	Entries []*blogs.BlogEntry
}

// Add links the given BlogEntry to this TagEntry.
func (te *TagEntry) Add(e *blogs.BlogEntry) {
	// If it needs to be initialized, do that now.
	if te.Entries == nil {
		te.Entries = make([]*blogs.BlogEntry, 0, 0)
	}

	te.Entries = append(te.Entries, e)
}
