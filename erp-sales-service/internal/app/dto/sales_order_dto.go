package dto

import "time"

type SalesOrderResponse struct {
	ID                  int64     `json:"id"`
	AddressIDGP         string    `json:"address_id_gp"`
	CustomerIDGP        string    `json:"customer_id_gp"`
	TermPaymentSlsIDGP  string    `json:"term_payment_sls_id_gp"`
	SubDistrictIDGP     string    `json:"sub_district_id_gp"`
	SiteIDGP            string    `json:"site_id_gp"`
	WrtIDGP             string    `json:"wrt_id_gp"`
	RegionIDGP          string    `json:"region_id_gp"`
	PriceLevelIDGP      string    `json:"price_level_id_gp"`
	PaymentGroupSlsID   int32     `json:"payment_group_sls_id"`
	ArchetypeIDGP       string    `json:"arechetype_id_gp"`
	SalesOrderNumber    string    `json:"sales_order_number"`
	IntegrationCode     string    `json:"integration_code"`
	SalesOrderNumberGP  string    `json:"sales_order_number_gp"`
	Status              int8      `json:"status"`
	RecognitionDate     time.Time `json:"recognition_date"`
	RequestsShipDate    time.Time `json:"requests_ship_date"`
	BillingAddress      string    `json:"billing_address"`
	ShippingAddress     string    `json:"shipping_address"`
	ShippingAddressNote string    `json:"shipping_address_note"`
	DeliveryFee         float64   `json:"delivery_fee"`
	VouDiscAmount       float64   `json:"vou_disc_amount"`
	CustomerPointLogID  int64     `json:"customer_point_log_id"`
	EdenPointCampaignID int64     `json:"den_point_campaign_id"`
	TotalPrice          float64   `json:"total_price"`
	TotalCharge         float64   `json:"total_charge"`
	TotalWeight         float64   `json:"total_weight"`
	Note                string    `json:"note"`
	PaymentReminder     int8      `json:"payment_reminder"`
	CancelType          int8      `json:"cancel_type"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           int64     `json:"created_by"`
	ShippingMethodIDGP  string    `json:"shipping_method_id_gp"`
	CustomerNameGP      string    `json:"customer_name_gp"`

	SalesOrderItem    []*SalesOrderItemResponse
	SalesOrderVoucher []*SalesOrderVoucherResponse
}

type SalesOrderFeedback struct {
	SalesOrderCode string  `json:"sales_order_code,omitempty"`
	DeliveryDate   string  `json:"delivery_date,omitempty"`
	SalesOrderID   int64   `json:"sales_order_id"`
	TotalCharge    float64 `json:"total_charge"`
	RatingScore    int8    `json:"rating_score"`
	Tags           string  `json:"tags,omitempty"`
	Description    string  `json:"description,omitempty"`
}
