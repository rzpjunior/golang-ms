package dto

import "time"

type ProspectiveCustomerResponse struct {
	ID               int64                 `json:"id"`
	Code             string                `json:"code"`
	Salesperson      *SalespersonResponse  `json:"salesperson"`
	Archetype        *ArchetypeResponse    `json:"archetype"`
	CustomerType     *CustomerTypeResponse `json:"customer_type"`
	SubDistrict      *SubDistrictResponse  `json:"sub_district"`
	Region           *RegionResponse       `json:"region"`
	Customer         *CustomerResponse     `json:"customer"`
	Name             string                `json:"name"`
	Phone1           string                `json:"phone_1"`
	Phone2           string                `json:"phone_2"`
	Phone3           string                `json:"phone_3"`
	CustomerUpgrade  int8                  `json:"customer_upgrade"`
	RegStatus        int8                  `json:"reg_status"`
	RegStatusConvert string                `json:"reg_status_convert"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	ProcessedAt      time.Time             `json:"processed_at"`
	ProcessedBy      int64                 `json:"processed_by"`
	DeclineType      int8                  `json:"decline_type"`
	DeclineNote      string                `json:"decline_note"`
}

type CustomerProspectiveResponse struct {
	ID   int64  `json:"-"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProspectiveCustomerDecineRequest struct {
	DeclineType int8   `json:"decline_type" valid:"required"`
	DeclineNote string `json:"decline_note" valid:"lte:250"`
}
