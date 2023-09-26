package dto

type ItemClassResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type GetItemClassRequest struct {
	Limit  int
	Offset int
	Search string
}
