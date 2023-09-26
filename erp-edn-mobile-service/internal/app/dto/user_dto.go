package dto

import "time"

type UserResponse struct {
	ID            int64           `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	Nickname      string          `json:"nickname,omitempty"`
	Email         string          `json:"email,omitempty"`
	Password      string          `json:"password,omitempty"`
	ParentID      int64           `json:"parent_id,omitempty"`
	SiteID        string          `json:"site_id,omitempty"`
	EmployeeCode  string          `json:"employee_code,omitempty"`
	PhoneNumber   string          `json:"phone_number,omitempty"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
	Status        int8            `json:"status,omitempty"`
	StatusConvert string          `json:"status_convert"`
	Region        *RegionResponse `json:"region,omitempty"`
}
