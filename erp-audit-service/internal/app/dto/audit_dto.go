package dto

import "time"

type AuditRequestCreate struct {
	UserID      int64  `json:"user_id"`
	UserIdGp    string `json:"user_id_gp"`
	ReferenceID string `json:"reference_id"`
	Type        string `json:"type"`
	Function    string `json:"function"`
	Note        string `json:"note"`
}

type AuditResponseCreate struct {
	ID          string    `json:"id"`
	UserID      int64     `json:"user_id"`
	UserIdGp    string    `json:"user_id_gp"`
	ReferenceID string    `json:"reference_id"`
	Type        string    `json:"type"`
	Function    string    `json:"function"`
	CreatedAt   time.Time `json:"created_at"`
	Note        string    `json:"note"`
}

type AuditResponseGet struct {
	ID          string        `json:"id" valid:"required"`
	UserID      int64         `json:"user_id" valid:"required"`
	UserIdGp    string        `json:"user_id_gp"`
	ReferenceID string        `json:"reference_id" valid:"required"`
	Type        string        `json:"type" valid:"required"`
	Function    string        `json:"function" valid:"required"`
	CreatedAt   time.Time     `json:"created_at" valid:"required"`
	Note        string        `json:"note"`
	User        *UserResponse `json:"user" valid:"required"`
}

type UserResponse struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	Email        string `json:"email,omitempty"`
	EmployeeCode string `json:"employee_code,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	MainRole     string `json:"main_role,omitempty"`
	Division     string `json:"division,omitempty"`
	Status       int8   `json:"status,omitempty"`
}
