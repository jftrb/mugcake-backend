package api

type PaginationRequest struct {
	Cursor string
}

type RecipeSummaryPaginationRequest struct {
	Limit  int
	Offset int
}