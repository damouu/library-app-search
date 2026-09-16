package domain

type SearchParams struct {
	Query string
	Page  int
	Size  int
}

type SearchResult struct {
	Items []Chapter `json:"items"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
}
