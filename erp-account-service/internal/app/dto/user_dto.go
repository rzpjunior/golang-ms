package dto

import (
	"time"
)

type UserResponse struct {
	ID                     int64              `json:"id,omitempty"`
	Name                   string             `json:"name,omitempty"`
	Nickname               string             `json:"nickname,omitempty"`
	Email                  string             `json:"email,omitempty"`
	Password               string             `json:"password,omitempty"`
	ParentID               int64              `json:"parent_id,omitempty"`
	Region                 *RegionResponse    `json:"region,omitempty"`
	Site                   *SiteResponse      `json:"site,omitempty"`
	Territory              *TerritoryResponse `json:"territory"`
	EmployeeCode           string             `json:"employee_code,omitempty"`
	PhoneNumber            string             `json:"phone_number,omitempty"`
	MainRole               *RoleResponse      `json:"main_role,omitempty"`
	SubRoles               []*RoleResponse    `json:"sub_roles,omitempty"`
	CreatedAt              time.Time          `json:"created_at,omitempty"`
	UpdatedAt              time.Time          `json:"updated_at,omitempty"`
	Status                 int8               `json:"status,omitempty"`
	StatusConvert          string             `json:"status_convert"`
	ForceLogout            int8               `json:"force_logout,omitempty"`
	Note                   string             `json:"note"`
	PurchaserAppLoginToken string             `json:"purchaser_login_token"`
	PurchaserAppNotifToken string             `json:"purchaser_notif_token"`
	Supervisor             *Supervisor        `json:"supervisor"`
}

type Supervisor struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	EmployeeCode string `json:"employee_code,omitempty"`
}
type UserRequestCreate struct {
	Name            string  `json:"name" valid:"required"`
	Nickname        string  `json:"nickname" valid:"required"`
	Email           string  `json:"email" valid:"required|email"`
	Password        string  `json:"password" valid:"required|gte:5"`
	PasswordConfirm string  `json:"password_confirm" valid:"required|gte:5"`
	RegionID        string  `json:"region_id" valid:"required"`
	ParentID        int64   `json:"parent_id"`
	SiteID          string  `json:"site_id" valid:"required"`
	TerritoryID     string  `json:"territory_id"`
	EmployeeCode    string  `json:"employee_code" valid:"required"`
	PhoneNumber     string  `json:"phone_number" valid:"required|numeric"`
	MainRole        int64   `json:"main_role" valid:"required"`
	SubRoles        []int64 `json:"sub_roles"`
	Note            string  `json:"note" valid:"gte=255"`
}

type UserRequestUpdate struct {
	Name        string  `json:"name" valid:"required"`
	Nickname    string  `json:"nickname" valid:"required"`
	RegionID    string  `json:"region_id" valid:"required"`
	ParentID    int64   `json:"parent_id"`
	SiteID      string  `json:"site_id" valid:"required"`
	TerritoryID string  `json:"territory_id"`
	PhoneNumber string  `json:"phone_number" valid:"required|numeric"`
	MainRole    int64   `json:"main_role" valid:"required"`
	SubRoles    []int64 `json:"sub_roles"`
	Note        string  `json:"note" valid:"gte=255"`
}

type UserRequestResetPassword struct {
	Password        string `json:"password" valid:"required|gte:5"`
	PasswordConfirm string `json:"confirm_password" valid:"required|gte:5"`
}

type GetUserRequest struct {
	Offset     int
	Limit      int
	Status     int
	Search     string
	OrderBy    string
	SiteID     string
	RegionID   string
	DivisionID int64
	RoleID     int64
	Apps       string
}
