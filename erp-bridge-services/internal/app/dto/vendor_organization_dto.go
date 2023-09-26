package dto

type VendorOrganizationResponse struct {
	ID                     int64  `json:"id"`
	Code                   string `json:"code"`
	VendorClassificationID int64  `json:"vendor_classification_id"`
	SubDistrictID          int64  `json:"sub_district_id"`
	PaymentTermID          int64  `json:"payment_term_id"`
	Name                   string `json:"name"`
	Address                string `json:"address"`
	Note                   string `json:"note"`
	Status                 int32  `json:"status"`
}
