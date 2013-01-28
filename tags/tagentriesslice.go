package tags

// TagEntriesSlice is a slice of TagEntry structures this is returns by
// the Slice() function for TagEntries. It implements the sorting
// interface for golangs sort package.
type TagEntriesSlice []*TagEntry

// Len returns the length of the TagEntriesSlice.
func (tes TagEntriesSlice) Len() int {
	return len(tes)
}

// Less returns true if the value at i is less than the value at j.
func (tes TagEntriesSlice) Less(i, j int) bool {
	return tes[i].Name < tes[j].Name
}

// Swap switches the elemens at i and j.
func (tes TagEntriesSlice) Swap(i, j int) {
	tes[i], tes[j] = tes[j], tes[i]
}
