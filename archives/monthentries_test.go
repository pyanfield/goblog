package archives

import (
	"github.com/icub3d/goblog/blogs"
	"sort"
	"testing"
	"time"
)

// entries is a list of blog entries that can be used in the testing
// environment. They should not be modified so that all the tests can
// use them.
var entries = []*blogs.BlogEntry{
	&blogs.BlogEntry{
		Name:    "Test Entry 1",
		Created: time.Unix(1999181, 0),
	},
	&blogs.BlogEntry{
		Name:    "Test Entry 2",
		Created: time.Unix(0, 0),
	},
	&blogs.BlogEntry{
		Name:    "Test Entry 3",
		Created: time.Unix(99381811, 0),
	},
	&blogs.BlogEntry{
		Name:    "Test Entry 4",
		Created: time.Unix(120, 0),
	},
	&blogs.BlogEntry{
		Name:    "Test Entry 5",
		Created: time.Unix(80, 0),
	},
}

// TestMonthEntriesAdd tests the MonthEntries.Add function.
func TestMonthEntriesAdd(t *testing.T) {
	// These are our test cases. We will iteratively add these and make
	// sure that the entries have been added.
	tests := []struct {
		entry *blogs.BlogEntry
	}{
		{
			entry: entries[0],
		},
		{
			entry: entries[1],
		},
		{
			entry: entries[2],
		},
	}

	// Make the MonthEntries.
	me := MonthEntries{
		Month: "November",
	}

	for i, test := range tests {
		// The original length of the entries.
		orglen := len(me.Entries)

		// Add the item.
		me.Add(test.entry)

		// Make sure we added the new item.
		if len(me.Entries) != orglen+1 {
			t.Errorf("(%d) Expecting %d from len(me.Entries) but got %d",
				i, orglen+1, len(me.Entries))
		}

		// Make sure the last one is the right one.
		if me.Entries[len(me.Entries)-1].Name != test.entry.Name {
			t.Errorf("(%d) Expecting %s from me.Entries[last].Name but got %s",
				i, test.entry.Name, me.Entries[len(me.Entries)-1].Name)
		}
	}
}

// TestMonthEntriesSort tests the sort functions (Len, Less, Swap) of
// MonthEntries by using sort.Sort.
func TestMonthEntriesSort(t *testing.T) {
	// These are our test cases.
	tests := []struct {
		me       *MonthEntries
		expected []*blogs.BlogEntry
	}{
		// A normal test.
		{
			me: &MonthEntries{
				Month: "November",
				Entries: []*blogs.BlogEntry{
					entries[0],
					entries[1],
					entries[2],
					entries[3],
					entries[4],
				},
			},
			expected: []*blogs.BlogEntry{
				entries[2],
				entries[0],
				entries[3],
				entries[4],
				entries[1],
			},
		},

		// An empty test.
		{
			me: &MonthEntries{
				Month:   "November",
				Entries: []*blogs.BlogEntry{},
			},
			expected: []*blogs.BlogEntry{},
		},
	}

	for i, test := range tests {
		// Sort the entries.
		sort.Sort(test.me)

		// Check the entries against the expected value.
		if len(test.expected) != len(test.me.Entries) {
			t.Errorf("(%d) expecting lengths of expected and test don't match: %d, %d",
				i, len(test.expected), len(test.me.Entries))
		}

		// Make sure each entry is in it's right place.
		for k, v := range test.expected {
			if test.me.Entries[k].Name != v.Name {
				t.Errorf("(%d) expecting '%s' from Entries[%d].name but got '%s'",
					i, v.Name, k, test.me.Entries[k].Name)
			}
		}
	}
}
