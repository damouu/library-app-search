package domain

type SearchParams struct {
	Query string
	Page  int
	Size  int
}

type SearchResult struct {
	Items      []Chapter
	Page       int
	Size       int
	Total      int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}
