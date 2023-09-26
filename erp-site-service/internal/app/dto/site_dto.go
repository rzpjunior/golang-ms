package dto

type SiteResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type FilterSiteResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get Site
type GetSiteRequest struct {
	Limit    int `json:"name"`
	Offset   int `json:"name"`
	Search string `json:"search"`
}
