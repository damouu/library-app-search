package domain

type Chapter struct {
	ChapterUUID     string `json:"chapter_uuid"`
	SeriesUUID      string `json:"series_uuid"`
	Title           string `json:"title"`
	SecondTitle     string `json:"second_title"`
	Summary         string `json:"summary"`
	ChapterNumber   int    `json:"chapter_number"`
	TotalPages      int    `json:"total_pages"`
	PublicationDate string `json:"publication_date"`
	CoverArtworkURL string `json:"cover_artwork_url"`
}
