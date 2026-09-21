package opensearch

import (
	"fmt"
	"library-app-search/internal/domain"
	"time"
	_ "time"
	"uuid"
)

func toDomainChapter(document ChapterDocument) (domain.Chapter, error) {
	chapterUUID, err := uuid.Parse(document.ChapterUUID)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("invalid chapter UUID: %w", err)
	}

	seriesUUID, err := uuid.Parse(document.SeriesUUID)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("invalid series UUID: %w", err)
	}

	publicationDate, err := time.Parse("2006-01-02", document.PublicationDate)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("invalid publication date: %w", err)
	}

	return domain.Chapter{
		ChapterUUID:     chapterUUID,
		SeriesUUID:      seriesUUID,
		Title:           document.Title,
		SecondTitle:     document.SecondTitle,
		Summary:         document.Summary,
		ChapterNumber:   document.ChapterNumber,
		TotalPages:      document.TotalPages,
		PublicationDate: publicationDate,
		CoverArtworkURL: document.CoverArtworkURL,
	}, nil
}
