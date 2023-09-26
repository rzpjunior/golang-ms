package dto

import "time"

type SalesPaymentTermResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SalesPaymentTermGPResponse struct {
	Pymtrmid              string  `json:"pymtrmid"`
	Duetype               int64   `json:"duetype"`
	Duedesc               string  `json:"duedesc"`
	Duedtds               int64   `json:"duedtds"`
	CalculateDateFrom     int64   `json:"calculateDateFrom"`
	CalculateDateFromDays int64   `json:"calculateDateFromDays"`
	Disctype              int64   `json:"disctype"`
	Discdtds              int64   `json:"discdtds"`
	Dsclctyp              int64   `json:"dsclctyp"`
	Dscdlram              float64 `json:"dscdlram"`
	Dscpctam              int64   `json:"dscpctam"`
	Salpurch              int64   `json:"salpurch"`
	Discntcb              int64   `json:"discntcb"`
	Freight               int64   `json:"freight"`
	Misc                  int64   `json:"misc"`
	Tax                   int64   `json:"tax"`
}
