package dto

import "time"

type VendorResponse struct {
	ID                     string                      `json:"id,omitempty"`
	Code                   string                      `json:"code,omitempty"`
	VendorOrganization     *VendorOrganizationResponse `json:"vendor_organization,omitempty"`
	VendorClassificationID int64                       `json:"vendor_classification_id,omitempty"`
	SubDistrictID          int64                       `json:"sub_district_id,omitempty"`
	PicName                string                      `json:"pic_name,omitempty"`
	Email                  string                      `json:"email,omitempty"`
	PhoneNumber            string                      `json:"phone_number,omitempty"`
	PaymentTermID          int64                       `json:"payment_term_id,omitempty"`
	Rejectable             int32                       `json:"rejectable,omitempty"`
	Returnable             int32                       `json:"returnable,omitempty"`
	Address                string                      `json:"address,omitempty"`
	Note                   string                      `json:"note,omitempty"`
	Status                 int32                       `json:"status,omitempty"`
	Latitude               string                      `json:"latitude,omitempty"`
	Longitude              string                      `json:"longitude,omitempty"`
	Name                   string                      `json:"name,omitempty"`
	CreatedAt              time.Time                   `json:"created_at,omitempty"`
	CreatedBy              int64                       `json:"created_by,omitempty"`
}

type VendorOrganizationResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type VendorListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type VendorDetailRequest struct {
	Id int32 `json:"id"`
}

type VendorGP struct {
	VendorId string `json:"id"`
	VendName string `json:"name"`
	Address  string `json:"address"`
	Inactive int32  `json:"inactive,omitempty"`
}

type GetVendorGPResponse struct {
	PageNumber   int32       `json:"pageNumber"`
	PageSize     int32       `json:"pageSize"`
	TotalPages   int32       `json:"totalPages"`
	TotalRecords int32       `json:"totalRecords"`
	Data         []*VendorGP `json:"data"`
	Succeeded    bool        `json:"succeeded"`
	Errors       []string    `json:"errors"`
	Message      string      `json:"message"`
}

type GetVendorGPListRequest struct {
	Limit   int32  `query:"limit"`
	Offset  int32  `query:"offset"`
	Status  int32  `query:"status"`
	Search  string `query:"search"`
	OrderBy string `query:"orderBy"`
}
