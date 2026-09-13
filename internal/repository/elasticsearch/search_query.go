package elasticsearch

type SearchQuery struct {
	Query QueryClause `json:"query"`
}

type QueryClause struct {
	MultiMatch MultiMatch `json:"multi_match"`
}

type MultiMatch struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}
