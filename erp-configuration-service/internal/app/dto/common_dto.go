package dto

type CommonFilterRequest struct {
	offset  int
	limit   int
	status  int
	search  string
	orderBy string
}
