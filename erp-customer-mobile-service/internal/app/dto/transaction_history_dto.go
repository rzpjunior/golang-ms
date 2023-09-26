package dto

import "time"

type RequestGetHistoryTransaction struct {
	Platform string                    `json:"platform" valid:"required"`
	Data     DataGetHistoryTransaction `json:"data" valid:"required"`
	Session  *SessionDataCustomer
	//Branch   *model.Branch
	Limit  int64 `json:"limit" valid:"required"`
	Offset int64 `json:"offset" `
}

type DataGetHistoryTransaction struct {
	AddressID  string `json:"address_id" valid:"required"`
	Type       string `json:"type" `
	Category   string `json:"category" `
	OrderBy    string `json:"order_by"`
	SalesOrder []*ListSalesOrder
}

type ListSalesOrder struct {
	ID                string `orm:"column(id);auto" json:"id"`
	OrderCode         string `orm:"column(code);size(50);null" json:"code"`
	OrderDate         string `orm:"column(recognition_date)" json:"recognition_date"`
	OrderDeliveryDate string `orm:"column(delivery_date)" json:"delivery_date"`
	OrderStatus       string `orm:"column(status)" json:"status"`
	TotalCharge       string `orm:"column(total_charge)" json:"total_charge"`
	OrderTypeSlsID    string `orm:"column(order_type_sls_id)" json:"order_type_sls_id"`
	TermPaymentSlsID  string `orm:"column(term_payment_sls_id)" json:"term_payment_sls_id"`
}

type RequestGetDetailSO struct {
	Platform string          `json:"platform" valid:"required"`
	Data     DataGetDetailSO `json:"data" valid:"required"`

	Session *SessionDataCustomer `json:"-"`
}

type DataGetDetailSO struct {
	AddressID    string `json:"address_id" valid:"required"`
	SalesOrderID string `json:"sales_order_id" valid:"required"`
}

type SalesOrderDetailResponse struct {
	ID                  string `orm:"-" json:"id"`
	OrderID             string `orm:"column(id)" json:"-"`
	Code                string `orm:"column(code);size(50);null" json:"code"`
	RecognitionDate     string `orm:"column(order_date)" json:"order_date"`
	DeliveryDate        string `orm:"column(delivery_date)" json:"order_date_delivery"`
	ShippingAddress     string `orm:"column(shipping_address)" json:"shipping_address"`
	ShippingAddressNote string `orm:"column(shipping_address_note)" json:"shipping_address_note"`
	DeliveryFee         string `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouDiscAmount       string `orm:"column(vou_disc_amount)" json:"vou_disc_amount"`
	VoucherType         string `orm:"column(voucher_type)" json:"voucher_type"`
	PointRedeemAmount   string `orm:"column(point_redeem_amount)" json:"point_redeem_amount"`
	TotalPrice          string `orm:"column(total_price)" json:"total_price"`
	TotalCharge         string `orm:"column(total_charge)" json:"total_charge"`
	Note                string `orm:"column(order_note)" json:"order_note"`
	Status              string `orm:"column(order_status)" json:"order_status"`
	CreatedAt           string `orm:"column(created_at)" json:"created_at"`
	WrtName             string `orm:"column(wrt_name)" json:"wrt_name"`

	SIID              string `orm:"-" json:"invoice_id"`
	InvoiceID         string `orm:"column(invoice_id)" json:"-"`
	InvoiceCode       string `orm:"column(invoice_code);size(50);null" json:"invoice_code"`
	InvoiceStatus     string `orm:"column(invoice_status)" json:"invoice_status"`
	HasExtInvoice     string `orm:"column(has_ext_invoice)" json:"has_ext_invoice"`
	PaymentGroupSlsId string `orm:"column(payment_group_sls_id)" json:"payment_group_sls_id"`

	AddressName string `orm:"column(address_name)" json:"address_name"`
	PicName     string `orm:"column(recipient_name)" json:"recipient_name"`
	PhoneNumber string `orm:"column(phone_number)" json:"phone_number"`

	CityName    string `orm:"column(city_name)" json:"city_name"`
	EditOrder   string `orm:"column(edit_order)" json:"edit_order"`
	CancelOrder string `orm:"column(cancel_order)" json:"cancel_order"`

	TermPaymentCode        string `orm:"column(tps_code)" json:"payment_term_code"`
	TermPaymentName        string `orm:"column(tps_name)" json:"payment_term_name"`
	TermPaymentDayValue    string `orm:"column(tps_day_value)" json:"payment_term_days_value"`
	TermPaymentNote        string `orm:"column(tps_note)" json:"payment_term_note"`
	TermPaymentDescription string `orm:"column(tps_description)" json:"payment_term_description"`
	TermPaymentStatus      string `orm:"column(tps_status)" json:"payment_term_status"`

	FinishedAt            time.Time `orm:"column(finished_at);type(timestamp);null" json:"finished_at"`
	IsHavePayment         bool      `json:"is_have_payment"`
	OrderTypeID           string    `orm:"column(order_type_sls_id)" json:"order_type_id"`
	TransactionStatus     string    `json:"transaction_status"`
	TransactionStatusName string    `json:"transaction_status_name"`

	SalesOrderItems []*SalesOrderItemResponse `json:"sales_order_items"`
}

type SalesOrderItemResponse struct {
	ID                 string `orm:"-" json:"id"`
	SalesOrderItemID   string `orm:"column(id)" json:"-"`
	ProductID          string `orm:"-" json:"item_id"`
	ProdID             string `orm:"column(product_id)" json:"-"`
	OrderQty           string `orm:"column(order_qty)" json:"order_qty"`
	UnitPrice          string `orm:"column(unit_price)" json:"unit_price"`
	ShadowPrice        string `orm:"column(shadow_price)" json:"shadow_price"`
	Subtotal           string `orm:"column(subtotal)" json:"subtotal"`
	Weight             string `orm:"column(unit_weight)" json:"weight"`
	Note               string `orm:"column(note)" json:"note"`
	UomName            string `orm:"column(uom_name)" json:"uom_name"`
	Name               string `orm:"column(name)" json:"name"`
	ImageUrl           string `orm:"column(image_url)" json:"image_url"`
	DiscountQty        string `orm:"column(discount_qty)" json:"discount_qty"`
	UnitPriceDiscount  string `orm:"column(unit_price_discount)" json:"unit_price_discount"`
	ItemDiscountAmount string `orm:"column(sku_disc_amount)" json:"item_discount_amount"`
}

type RequestGetInvoiceDetail struct {
	Platform string           `json:"platform" valid:"required"`
	Data     GetDetailInvoice `json:"data" valid:"required"`

	Session *SessionDataCustomer `json:"-"`
}

type GetDetailInvoice struct {
	AddressID      string `json:"address_id" valid:"required"`
	SalesOrderID   string `json:"sales_order_id" valid:"required"`
	SalesInvoiceID string `json:"sales_invoice_id" valid:"required"`
}

type SalesInvoiceDetailResponse struct {
	InvoiceDetail  *SalesInvoice
	InvoicePayment []*SalesPayment
	InvoiceItem    []*SalesInvoiceItem
}

type SalesInvoice struct {
	ID                string  `json:"id"`
	InvoiceID         int64   `json:"-"`
	OrderCode         string  `json:"order_code"`
	InvoiceCode       string  `json:"invoice_code"`
	OrderDate         string  `json:"order_date"`
	InvoiceDate       string  `json:"invoice_date"`
	TotalPrice        float64 `json:"total_price"`
	DeliveryFee       float64 `json:"delivery_fee"`
	VoucherAmount     float64 `json:"voucher_amount"`
	PointRedeemAmount float64 `json:"point_redeem_amount"`
	AdjustmentAmount  float64 `json:"adjustment_amount"`
	TotalCharge       float64 `json:"total_charge"`
}

type SalesPayment struct {
	ID             string `json:"sales_payment_id"`
	Code           string `json:"sales_payment_code"`
	PaymentDate    string `json:"payment_date"`
	PaymentMethod  string `json:"payment_method"`
	PaymentTime    string `json:"payment_time"`
	Amount         string `json:"amount"`
	PaymentChannel string `json:"payment_channel"`
	Status         int8   `json:"status"`
}

type SalesInvoiceItem struct {
	ItemID            string  `json:"item_id"`
	ItemName          string  `json:"item_name"`
	InvoiceQty        float64 `json:"invoice_qty"`
	UomName           string  `json:"uom_name"`
	UnitPrice         float64 `json:"unit_price"`
	Subtotal          float64 `json:"subtotal"`
	SkuDiscountAmount float64 `json:"sku_disc_amount"`
}

type RequestGetDetailSOInvoice struct {
	Platform string                 `json:"platform" valid:"required"`
	Data     dataGetDetailSOInvoice `json:"data" valid:"required"`
	Session  *SessionDataCustomer
	// SalesOrder   *model.SalesOrder
	// SalesInvoice *model.SalesInvoice
}

type dataGetDetailSOInvoice struct {
	BranchID  string `json:"branch_id" valid:"required"`
	SoID      string `json:"sales_order_id" valid:"required"`
	InvoiceID string `json:"invoice_id" valid:"required"`

	// DataResponse *model.SalesInvoiceResponse `json:"-"`
	// Branch       *model.Branch               `json:"-"`
}

type ResponseDetailSOInvoice struct {
	InvoiceDetail  *SalesInvoice
	InvoicePayment []*ListInvoicePayment
	InvoiceItem    []*ListInvoiceItem
}

// ListInvoicePayment : struct to hold list of sales payment data
type ListInvoicePayment struct {
	ID             string `orm:"-" json:"sales_payment_id"`
	PaymentID      int64  `orm:"column(sales_payment_id)" json:"-"`
	Code           string `orm:"column(sales_payment_code);size(50);null" json:"sales_payment_code"`
	PaymentDate    string `orm:"column(payment_date);type(timestamp);null" json:"payment_date"`
	PaymentMethod  string `orm:"column(payment_method)" json:"payment_method"`
	PaymentTime    string `orm:"column(payment_time);type(timestamp);null" json:"payment_time"`
	Amount         string `orm:"column(amount)" json:"amount"`
	PaymentChannel string `orm:"column(payment_channel)" json:"payment_channel"`
	Status         int8   `orm:"column(status)" json:"status"`
}

// ListInvoiceItem : struct to hold list of sales invoice item data
type ListInvoiceItem struct {
	ProductIDInt      int64   `orm:"column(product_id)" json:"-"`
	ProductID         string  `orm:"-" json:"product_id"`
	ProductName       string  `orm:"column(product_name)" json:"product_name"`
	InvoiceQty        float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	UomName           string  `orm:"column(uom_name)" json:"uom_name"`
	UnitPrice         float64 `orm:"column(unit_price)" json:"unit_price"`
	Subtotal          float64 `orm:"column(subtotal)" json:"subtotal"`
	SkuDiscountAmount float64 `orm:"column(sku_disc_amount)" json:"sku_disc_amount"`
}
