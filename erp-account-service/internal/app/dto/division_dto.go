package dto

import "time"

type DivisionResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Status        int8      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StatusConvert string    `json:"status_convert"`
	Note          string    `json:"note"`
}

type DivisionRequestCreate struct {
	Code string `json:"code" valid:"required"`
	Name string `json:"name" valid:"required"`
	Note string `json:"note" valid:"lte:250"`
}

type DivisionRequestUpdate struct {
	Code string `json:"code" valid:"required"`
	Name string `json:"name" valid:"required"`
	Note string `json:"note" valid:"lte:250"`
}
