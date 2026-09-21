package handler

import (
	"library-app-search/internal/domain"
)

func toChapterResponse(chapter domain.Chapter) ChapterResponse {
	return ChapterResponse{
		ChapterUUID:     chapter.ChapterUUID.String(),
		SeriesUUID:      chapter.SeriesUUID.String(),
		Title:           chapter.Title,
		SecondTitle:     chapter.SecondTitle,
		Summary:         chapter.Summary,
		ChapterNumber:   chapter.ChapterNumber,
		TotalPages:      chapter.TotalPages,
		PublicationDate: chapter.PublicationDate.Format("2006-01-02"),
		CoverArtworkURL: chapter.CoverArtworkURL,
	}
}
