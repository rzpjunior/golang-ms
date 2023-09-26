package dto

import "time"

type SalesOrderResponse struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	DocNumber     string    `json:"doc_number"`
	AddressID     string    `json:"address_id"`
	CustomerID    string    `json:"customer_id"`
	SalespersonID string    `json:"salesperson_id"`
	WrtID         string    `json:"wrt_id"`
	PaymentTermID string    `json:"payment_term_id"`
	OrderTypeID   string    `json:"order_type_id"`
	SiteID        string    `json:"site_id"`
	Application   int8      `json:"application"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	OrderDate     time.Time `json:"order_date"`
	Total         float64   `json:"total"`
	CreatedDate   time.Time `json:"created_date"`
	ModifiedDate  time.Time `json:"modified_date"`
	FinishedDate  time.Time `json:"finished_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	AddressGP *AddressResponse `json:"address"`
	Site      *SiteResponse    `json:"site"`
}

type CreateSalesOrderRequest struct {
	CustomerID         string `json:"customer_id" valid:"required"`
	AddressID          string `json:"address_id" valid:"required"`
	RegionID           string `json:"region_id" valid:"required"`
	DeliveryDateStr    string `json:"delivery_date" valid:"required"`
	WrtID              string `json:"wrt_id" valid:"required"`
	SiteID             string `json:"site_id" valid:"required"`
	RecognitionDateStr string `json:"order_date" valid:"required"`
	OrderTypeID        string `json:"order_type_id" valid:"required"`
	SalespersonID      string `json:"salesperson_id" valid:"required"`
	SalesTermID        string `json:"term_payment_sls_id" valid:"required"`
	InvoiceTermID      string `json:"term_invoice_sls_id" valid:"required"`
	PaymentGroupID     string `json:"payment_group_id" valid:"required"`
	ShippingAddress    string `json:"shipping_address" valid:"required"`
	BillingAddress     string `json:"billing_address" valid:"required"`
	RedeemCode         string `json:"redeem_code"`
	Note               string `json:"note"`

	Items []*salesOrderItem `json:"items" valid:"required"`
}

type salesOrderItem struct {
	ItemID      string  `json:"item_id"`
	Quantity    float64 `json:"qty"`
	UnitPrice   float64 `json:"unit_price"`
	Note        string  `json:"note"`
	ProductPush int8    `json:"product_push"`
}
