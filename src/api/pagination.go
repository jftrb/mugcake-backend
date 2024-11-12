package api

type PaginationRequest struct {
	Cursor string
	Limit int
}

type RecipeSummaryPaginationRequest struct {
	Offset int
}