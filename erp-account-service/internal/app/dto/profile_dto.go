package dto

import "time"

type ProfileResponse struct {
	ID            int64           `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	Nickname      string          `json:"nickname,omitempty"`
	Email         string          `json:"email,omitempty"`
	Password      string          `json:"password,omitempty"`
	RegionID      int64           `json:"region_id,omitempty"`
	ParentID      int64           `json:"parent_id,omitempty"`
	SiteID        int64           `json:"site_id,omitempty"`
	TerritoryID   int64           `json:"territory_id,omitempty"`
	EmployeeCode  string          `json:"employee_code,omitempty"`
	PhoneNumber   string          `json:"phone_number,omitempty"`
	Region        *RegionResponse `json:"region,omitempty"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
	Status        int8            `json:"status,omitempty"`
	StatusConvert string          `json:"status_convert"`
}

type ProfileRequestUpdate struct {
	Nickname    string `json:"nickname" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required|numeric"`
}

type UpdateSalesAppTokenRequest struct {
	Id                 int64  `json:"id" valid:"required"`
	SalesAppLoginToken string `json:"salesapp_login_token" valid:"required"`
	SalesAppNotifToken string `json:"salesapp_notif_token"`
	ForceLogout        int32  `json:"force_logout"`
}

type UpdateEdnAppTokenRequest struct {
	Id               int64  `json:"id" valid:"required"`
	EdnAppLoginToken string `json:"ednapp_login_token" valid:"required"`
	ForceLogout      int32  `json:"force_logout"`
}

type UpdatePurchaserAppTokenRequest struct {
	Id                     int64  `json:"id" valid:"required"`
	PurchaserAppNotifToken string `json:"purchaserapp_notif_token" valid:"required"`
	ForceLogout            int32  `json:"force_logout"`
}
