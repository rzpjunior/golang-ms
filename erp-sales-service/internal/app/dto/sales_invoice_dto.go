package dto

type ResponseSIMobile struct {
	InvoiceDetail  *SalesInvoice
	InvoicePayment []*ListInvoicePayment
	InvoiceItem    []*ListInvoiceItem
}

type SalesInvoice struct {
	ID                string  `json:"id"`
	InvoiceID         string  `json:"-"`
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
