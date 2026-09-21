package handler

import "library-app-search/internal/domain"

func toSearchResponse(result domain.SearchResult) SearchResponse {
	items := make([]ChapterResponse, 0, len(result.Items))

	for _, chapter := range result.Items {
		items = append(items, toChapterResponse(chapter))
	}

	return SearchResponse{
		Items:      items,
		Page:       result.Page,
		Size:       result.Size,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}
}
