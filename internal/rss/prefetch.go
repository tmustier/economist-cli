package rss

// PrefetchAll fetches all known sections to warm the RSS cache.
func PrefetchAll() {
	sections := SectionList()
	for _, info := range sections {
		_, _ = FetchSection(info.Path)
	}
}
