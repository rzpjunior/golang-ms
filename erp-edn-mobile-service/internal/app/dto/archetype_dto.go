package dto

import "time"

type ArchetypeResponse struct {
	ID             int64                 `json:"id"`
	Code           string                `json:"code"`
	CustomerTypeID int64                 `json:"customer_type_id"`
	Description    string                `json:"description"`
	Status         int8                  `json:"status"`
	StatusConvert  string                `json:"status_convert"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	CustomerType   *CustomerTypeResponse `json:"customer_type"`
}

type ArchetypeListRequest struct {
	Limit          int32  `json:"limit"`
	Offset         int32  `json:"offset"`
	Status         int32  `json:"status"`
	Search         string `json:"search"`
	OrderBy        string `json:"order_by"`
	CustomerTypeID int64  `json:"customer_type_id"`
}

type ArchetypeDetailRequest struct {
	Id int32 `json:"id"`
}

type GetArchetypeGPListRequest struct {
	Limit                   int32  `json:"limit"`
	Offset                  int32  `json:"offset"`
	GnlArchetypeId          string `json:"gnl_archetype_id"`
	GnlArchetypedescription string `json:"gnl_archetypedescription"`
	GnlCustTypeId           string `json:"gnl_cust_type_id"`
	Inactive                string `json:"inactive"`
}

type ArchetypeGP struct {
	GnlArchetypeId          string `json:"gnl_archetype_id"`
	GnlArchetypedescription string `json:"gnl_archetypedescription"`
	GnlCustTypeId           string `json:"gnl_cust_type_id"`
	GnlCusttypeDescription  string `json:"gnl_custtype_description"`
	Inactive                int32  `json:"inactive"`
	InactiveDesc            string `json:"inactive_desc"`
}
