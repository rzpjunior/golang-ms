package dto

type DeliveryDateResponse struct {
	Date []string
}

type DeliveryDateRequest struct {
	Platform string      `json:"platform" valid:"required"`
	Data     dataGetDate `json:"data" valid:"required"`
}

type dataGetDate struct {
	RegionID     string `json:"region_id" valid:"required"`
	DataResponse dataResponse
}

type dataResponse struct {
	Date []string
}
