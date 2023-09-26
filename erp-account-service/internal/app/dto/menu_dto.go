package dto

type MenuResponse struct {
	ID       int64           `json:"id,omitempty"`
	ParentID int64           `json:"parent_id,omitempty"`
	Title    string          `json:"title,omitempty"`
	Url      string          `json:"url"`
	Icon     string          `json:"icon,omitempty"`
	Child    []*MenuResponse `json:"child,omitempty"`
}
