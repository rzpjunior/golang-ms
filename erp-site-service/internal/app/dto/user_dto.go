package dto

import "time"

type UserResponse struct {
	ID                 int64     `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Nickname           string    `json:"nickname,omitempty"`
	Email              string    `json:"email,omitempty"`
	Password           string    `json:"password,omitempty"`
	RegionID           int64     `json:"region_id,omitempty"`
	ParentID           int64     `json:"parent_id,omitempty"`
	SiteID             int64     `json:"site_id,omitempty"`
	TerritoryID        int64     `json:"territory_id,omitempty"`
	EmployeeCode       string    `json:"employee_code,omitempty"`
	PhoneNumber        string    `json:"phone_number,omitempty"`
	Status             int8      `json:"status,omitempty"`
	Note               string    `json:"note,omitempty"`
	ForceLogout        int32     `json:"force_logout,omitempty"`
	SalesAppLoginToken string    `json:"salesapp_login_token,omitempty"`
	SalesAppNotifToken string    `json:"salesapp_notif_token,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

type GetUserRequest struct {
	Offset     int
	Limit      int
	Status     int
	Search     string
	OrderBy    string
	SiteId     int64
	DivisionId int64
	RoleId     int64
}
