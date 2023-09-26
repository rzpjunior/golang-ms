package dto

import "time"

type FieldPurchaserResponse struct {
	ID            int64              `json:"id"`
	Name          string             `json:"name"`
	Nickname      string             `json:"nickname,omitempty"`
	Email         string             `json:"email,omitempty"`
	Password      string             `json:"password,omitempty"`
	ParentID      int64              `json:"parent_id,omitempty"`
	Region        *RegionResponse    `json:"region,omitempty"`
	Site          *SiteResponse      `json:"site,omitempty"`
	Territory     *TerritoryResponse `json:"territory,omitempty"`
	EmployeeCode  string             `json:"employee_code,omitempty"`
	PhoneNumber   string             `json:"phone_number,omitempty"`
	Division      string             `json:"division,omitempty"`
	MainRole      string             `json:"main_role,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty"`
	Status        int8               `json:"status,omitempty"`
	StatusConvert string             `json:"status_convert,omitempty"`
	ForceLogout   int8               `json:"force_logout,omitempty"`
	Note          string             `json:"note,omitempty"`
}

type FieldPurchaserListRequest struct {
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
	Status   int32  `json:"status"`
	Search   string `json:"search"`
	OrderBy  string `json:"order_by"`
	SiteID   int32  `json:"site_id"`
	SiteIDGp string `json:"-"`
}
