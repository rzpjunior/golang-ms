package dto

import "time"

type SiteResponse struct {
	ID             string          `json:"id,omitempty"`
	Code           string          `json:"code,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Status         int8            `json:"status,omitempty"`
	StatusConvert  string          `json:"status_convert,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at,omitempty"`
	Address        string          `json:"address,omitempty"`
	Region         *RegionResponse `json:"region,omitempty"`
	PhoneNumber    string          `json:"phone_number,omitempty"`
	AltPhoneNumber string          `json:"alt_phone_number,omitempty"`
}

type SiteListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type SiteDetailRequest struct {
	Id int32 `json:"id"`
}

type SiteGP struct {
	ID                      string `json:"id,omitempty"`
	GnL_Site_Type_ID        string `json:"gnL_Site_Type_ID,omitempty"`
	Inactive                int32  `json:"inactive,omitempty"`
	Name                    string `json:"name,omitempty"`
	Address                 string `json:"address,omitempty"`
	AddresS2                string `json:"addresS2,omitempty"`
	AddresS3                string `json:"addresS3,omitempty"`
	PhonE1                  string `json:"phonE1,omitempty"`
	PhonE2                  string `json:"phonE2,omitempty"`
	PhonE3                  string `json:"phonE3,omitempty"`
	City                    string `json:"city,omitempty"`
	State                   string `json:"state,omitempty"`
	Faxnumbr                string `json:"faxnumbr,omitempty"`
	Zipcode                 string `json:"zipcode,omitempty"`
	Ccode                   string `json:"ccode,omitempty"`
	Country                 string `json:"country,omitempty"`
	Staxschd                string `json:"staxschd,omitempty"`
	Pctaxsch                string `json:"pctaxsch,omitempty"`
	Location_Segment        string `json:"location_Segment,omitempty"`
	GnL_Administrative_Code string `json:"gnL_Administrative_Code,omitempty"`
	Description             string `json:"description"`
}
