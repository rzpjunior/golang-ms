package dto

import "time"

type CreateRequestSalesOrder struct {
	Platform string         `json:"platform" valid:"required"`
	Data     dataSalesOrder `json:"data" valid:"required"`
	Session  *SessionDataCustomer
}

type dataSalesOrder struct {
	Code                      string    `json:"-"`
	AddressID                 string    `json:"address_id" valid:"required"`
	RegionID                  string    `json:"region_id" valid:"required"`
	DeliveryDateStr           string    `json:"delivery_date" valid:"required"`
	WrtID                     string    `json:"wrt_id" valid:"required"`
	BillingAddress            string    `json:"billing_address" valid:"required"`
	RedeemCode                string    `json:"redeem_code"`
	Note                      string    `json:"note"`
	RedeemPoint               float64   `json:"redeem_point"`
	TotalPrice                float64   `json:"-"`
	TotalWeight               float64   `json:"-"`
	TotalCharge               float64   `json:"-"`
	DeliveryFee               float64   `json:"-"`
	SameTagCustomer           string    `json:"-"`
	DeliveryDate              time.Time `json:"-"`
	RecognitionDate           time.Time `json:"-"`
	TotalSkuDiscount          float64   `json:"-"`
	IntegrationCode           string    `json:"-"`
	IsInitCustomerTalonPoints bool      `json:"-"`
	IsInitReferrerTalonPoints bool      `json:"-"`
	ReferrerData              []string  `json:"-"`
	CreditLimitBefore         float64   `json:"-"`
	CreditLimitAfter          float64   `json:"-"`
	IsCreateCreditLimitLog    bool      `json:"-"`
	ShippingAddress           string    `json:"shipping_address"`
	OrderChannel              int8      `json:"-"`
	OrderTypeID               int64     `json:"order_type_id"`
	CurrentPeriodPoint        float64   `json:"-"`
	NextPeriodPoint           float64   `json:"-"`
	CurrentPointUsed          float64   `json:"-"`
	NextPointUsed             float64   `json:"-"`

	IsCreateMerchantVa map[string]int8 `json:"-"`
	Second             float64         `json:"-"`

	Items       []*salesOrderItem `json:"items" valid:"required"`
	RecentPoint float64           `json:"-"`
	Payment     *payment          `json:"payment"`
}

type payment struct {
	PaymentMethod  string `json:"payment_method"`
	PaymentChannel string `json:"payment_channel"`
}

type salesOrderItem struct {
	ItemID            string  `json:"item_id"`
	Quantity          float64 `json:"qty"`
	UnitPrice         float64 `json:"unit_price"`
	Note              string  `json:"note"`
	ItemtPush         int8    `json:"item_push"`
	TaxableItem       int8    `json:"-"`
	TaxPercentage     float64 `json:"-"`
	DiscQty           float64 `json:"disc_qty"`
	UnitPriceDiscount float64 `json:"-"`
	DiscAmount        float64 `json:"-"`
	Subtotal          float64 `json:"-"`
	Weight            float64 `json:"-"`
	IsUseSkuDiscount  int8    `json:"-"`
}

// createRequest : struct to hold sales order request data
type UpdateCodRequest struct {
	Platform string     `json:"platform" valid:"required"`
	Data     dataUpdate `json:"data" valid:"required"`
	Session  *SessionDataCustomer
}

type dataUpdate struct {
	SalesOrderID string `json:"sales_order_id" valid:"required"`
}

type GetFeedback struct {
	Platform string           `json:"platform" valid:"required"`
	Offset   int64            `json:"offset"`
	Limit    int64            `json:"limit"`
	Data     dataTypeFeedback `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataTypeFeedback struct {
	FeedbackType int64 `json:"feedback_type"`

	DataResponse []dataResponse
}

type SalesOrderFeedback struct {
	ID             int64     ` json:"-"`
	SalesOrderCode string    `json:"sales_order_code,omitempty"`
	DeliveryDate   string    `json:"delivery_date,omitempty"`
	RatingScore    string    `json:"rating_score"`
	Tags           string    `json:"tags,omitempty"`
	Description    string    `json:"description,omitempty"`
	ToBeContacted  int8      `json:"-"`
	CreatedAt      time.Time `json:"-"`
	TotalCharge    string    `json:"total_charge"`
	SalesOrder     string    `json:"sales_order_id,omitempty"`
	// Merchant       *Merchant   `json:"-"`
}

type CreateSalesFeedback struct {
	Platform string       `json:"platform" valid:"required"`
	Data     dataFeedback `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataFeedback struct {
	SalesOrderId   string   `json:"sales_order_id"`
	SalesOrderCode string   `json:"sales_order_code"`
	DeliveryDate   string   `json:"delivery_date"`
	RatingScore    int      `json:"rating_score"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description"`
	ToBeContacted  int      `json:"to_be_contacted"`
	ExistTags      string
	CustomerId     int64
	//SalesOrder     *model.SalesOrder
}
