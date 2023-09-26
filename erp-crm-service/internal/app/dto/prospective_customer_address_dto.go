package dto

import "time"

type ProspectiveCustomerAddressResponse struct {
	ID                    int64     `json:"id"`
	ProspectiveCustomerID int64     `json:"prospective_customer_id"`
	AdmDivisionID         string    `json:"adm_division_id"`
	AddressName           string    `json:"address_name"`
	AddressType           string    `json:"address_type"`
	Address1              string    `json:"address_1"`
	Address2              string    `json:"address_2"`
	Address3              string    `json:"address_3"`
	Region                string    `json:"region"`
	Province              string    `json:"province"`
	City                  string    `json:"city"`
	District              string    `json:"district"`
	SubDistrict           string    `json:"sub_district"`
	PostalCode            string    `json:"postal_code"`
	Note                  string    `json:"note"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	ReferTo               int8      `json:"refer_to"`
}

type ProspectiveCustomerAddressGetDetailRequest struct {
	ID                    int64  `json:"id"`
	ProspectiveCustomerID int64  `json:"prospective_customer_id"`
	AddressType           string `json:"address_type"`
}
