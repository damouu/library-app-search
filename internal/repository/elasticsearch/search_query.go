package elasticsearch

type SearchQuery struct {
	From  int         `json:"from"`
	Size  int         `json:"size"`
	Query QueryClause `json:"query"`
}

type QueryClause struct {
	MultiMatch MultiMatch `json:"multi_match"`
}

type MultiMatch struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}
