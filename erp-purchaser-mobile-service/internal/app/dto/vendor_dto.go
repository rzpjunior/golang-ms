package dto

import "time"

type VendorResponse struct {
	ID                   string                        `json:"id"`
	Code                 string                        `json:"code"`
	VendorOrganization   *VendorOrganizationResponse   `json:"vendor_organization"`
	VendorClassification *VendorClassificationResponse `json:"vendor_classification"`
	AdmDivision          *AdmDivisionResponse          `json:"adm_division"`
	Name                 string                        `json:"name"`
	PicName              string                        `json:"pic_name"`
	Email                string                        `json:"email"`
	PhoneNumber          string                        `json:"phone_number"`
	PhoneNumberAlt       string                        `json:"phone_number_alt"`
	PaymentTerm          *PaymentTermResponse          `json:"payment_term"`
	PaymentMethod        *PaymentMethodResponse        `json:"payment_method"`
	Rejectable           int32                         `json:"rejectable"`
	Returnable           int32                         `json:"returnable"`
	Address              string                        `json:"address"`
	Note                 string                        `json:"note"`
	Status               int32                         `json:"status"`
	Latitude             string                        `json:"latitude"`
	Longitude            string                        `json:"longitude"`
	CreatedAt            time.Time                     `json:"created_at"`
	CreatedBy            int64                         `json:"created_by"`
}

type VendorListRequest struct {
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
	Status    int32  `json:"status"`
	Search    string `json:"search"`
	OrderBy   string `json:"orderby"`
	VendorOrg string `json:"vendor_org"`
	State     string `json:"state"`
}

type VendorRequestCreate struct {
	Name                   string `json:"name" valid:"required"`
	PicName                string `json:"pic_name" valid:"required"`
	PhoneNumber            string `json:"phone_number" valid:"required|numeric|range:8,15"`
	AltPhoneNumber         string `json:"alt_phone_number" valid:"numeric"`
	PaymentMethodID        string `json:"payment_method_id" valid:"required"`
	Rejectable             int8   `json:"rejectable" valid:"required"`
	Returnable             int8   `json:"returnable" valid:"required"`
	Address                string `json:"address" valid:"required"`
	BlockNumber            string `json:"block_number" valid:"lte:10"`
	VendorOrganizationID   string `json:"vendor_organization_id"`
	VendorClassificationID string `json:"vendor_classification_id"`
	PaymentTermID          string `json:"payment_term_id"`
}

type VendorRequestCreateResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
