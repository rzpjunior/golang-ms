package dto

type GetHelperResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetHelperRequest struct {
	Offset int
	Limit  int
	SiteId string
	Role   string
	Name   string
	Type   string
}
