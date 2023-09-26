package dto

type CourierVendorResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	SiteId string `json:"site_id"`
	Status int32  `json:"status"`
}

// Get Courier Vendor
type GetCourierVendorRequest struct {
	Limit             int
	Offset            int
	SiteId            string
	CourierVendorName string
}
