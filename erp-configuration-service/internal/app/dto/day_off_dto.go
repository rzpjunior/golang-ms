package dto

import "time"

type DayOffResponse struct {
	ID            int64     `json:"id,omitempty"`
	OffDate       time.Time `json:"off_date"`
	Note          string    `json:"note"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
}

type DayOffRequestCreate struct {
	OffDate time.Time `json:"off_date"`
	Note    string    `json:"note"`
}

type DayOffRequestUpdate struct {
	OffDate time.Time `json:"off_date"`
	Note    string    `json:"note"`
}
