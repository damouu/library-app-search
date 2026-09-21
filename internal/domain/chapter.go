package domain

import (
	"time"
	"uuid"
)

type Chapter struct {
	ChapterUUID     uuid.UUID
	SeriesUUID      uuid.UUID
	Title           string
	SecondTitle     string
	Summary         string
	ChapterNumber   int
	TotalPages      int
	PublicationDate time.Time
	CoverArtworkURL string
}
