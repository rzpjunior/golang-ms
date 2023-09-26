package dto

type UomGPResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type GetUomRequest struct {
	Limit  int
	Offset int
	Search string
}
