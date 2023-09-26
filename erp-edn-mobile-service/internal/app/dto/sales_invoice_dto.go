package dto

import (
	"time"
)

type SalesInvoiceResponse struct {
	ID                 int64     `json:"-"`
	Code               string    `json:"code"`
	CodeExt            string    `json:"code_ext"`
	Status             int8      `json:"status"`
	RecognitionDate    time.Time `json:"recognition_date"`
	DueDate            time.Time `json:"due_date"`
	BillingAddress     string    `json:"billing_address"`
	DeliveryFee        float64   `json:"delivery_fee"`
	VouRedeemCode      string    `json:"vou_redeem_code"`
	VouDiscAmount      float64   `json:"vou_disc_amount"`
	PointRedeemAmount  float64   `json:"point_redeem_amount"`
	Adjustment         int8      `json:"adjustment"`
	AdjAmount          float64   `json:"adj_amount"`
	AdjNote            string    `json:"adj_note"`
	TotalPrice         float64   `json:"total_price"`
	TotalCharge        float64   `json:"total_charge"`
	DeltaPrint         int64     `json:"delta_print"`
	Note               string    `json:"note"`
	VoucherID          int64     `json:"voucher_id"`
	RemainingAmount    float64   `json:"remaining_amount"`
	TotalPaid          float64   `json:"total_paid"`
	SiteID             string    `json:"site_id"`
	CustomerTypeID     string    `json:"customer_type_id"`
	CustomerTypeDesc   string    `json:"customer_type_desc"`
	CustomerName       string    `json:"customer_name"`
	RegionID           string    `json:"region_id"`
	CreatedAt          time.Time `json:"created_at"`
	VoucherType        int8      `orm:"-" json:"voucher_type"`
	TotalSkuDiscAmount float64   `orm:"column(total_sku_disc_amount)" json:"total_sku_disc_amount"`

	SalesOrder       *SalesOrderResponse         `json:"sales_order"`
	SalesInvoiceItem []*SalesInvoiceItemResponse `json:"sales_invoice_item"`
	SalesPayment     []*SalesPayment             `json:"sales_payment,omitempty"`
	// SalesInvoice *SalesInvoice
}

type SalesInvoiceListRequest struct {
	Limit                     int32     `json:"limit"`
	Offset                    int32     `json:"offset"`
	Status                    string    `json:"status"`
	Search                    string    `json:"search"`
	OrderBy                   string    `json:"order_by"`
	SiteID                    string    `json:"-"`
	RecognitionDateFrom       time.Time `json:"-"`
	RecognitionDateTo         time.Time `json:"-"`
	RecognitionDateFromString string    `json:"recognition_date_from"`
	RecognitionDateToString   string    `json:"recognition_date_to"`
	CustomerID                string    `json:"-"`
	RemainingAmountFlag int8 `json:"remaining_amount_flag"`
}

type SalesInvoiceDetailRequest struct {
	Id string `json:"id"`
}

type VoucherApply struct {
	GnlVoucherType int32   `json:"gnl_voucher_type"`
	GnlVoucherId   string  `json:"gnl_voucher_id"`
	Ordocamt       float64 `json:"ordocamt"`
}

type AmountReceived struct {
	Amount   float64 `json:"amount"`
	Chekbkid string  `json:"chekbkid"`
}

type DetailItem struct {
	Lnitmseq  int32   `json:"lnitmseq"`
	Itemnmbr  string  `json:"itemnmbr"`
	Locncode  string  `json:"locncode"`
	Uofm      string  `json:"uofm"`
	Pricelvl  string  `json:"pricelvl"`
	Quantity  int32   `json:"quantity"`
	Unitprce  float64 `json:"unitprce"`
	Xtndprce  float64 `json:"xtndprce"`
	GnlWeight int32   `json:"gnL_Weight"`
}

type CreateSalesInvoiceGPRequest struct {
	Interid            string          `json:"interid"`
	Orignumb           string          `json:"orignumb"`
	Sopnumbe           string          `json:"sopnumbe"`
	Docid              string          `json:"docid"`
	Docdate            string          `json:"docdate"`
	Custnmbr           string          `json:"custnmbr"`
	Custname           string          `json:"custname"`
	Prstadcd           string          `json:"prstadcd"`
	Curncyid           string          `json:"curncyid"`
	Subtotal           float64         `json:"subtotal"`
	Trdisamt           float64         `json:"trdisamt"`
	Freight            float64         `json:"freight"`
	Miscamnt           float64         `json:"miscamnt"`
	Taxamnt            float64         `json:"taxamnt"`
	Docamnt            float64         `json:"docamnt"`
	GnlRequestShipDate string          `json:"gnl_request_ship_date"`
	GnlRegion          string          `json:"gnl_region"`
	GnlWrtID           string          `json:"gnl_wrt_id"`
	GnlArchetypeID     string          `json:"gnl_archetype_id"`
	GnlOrderChannel    string          `json:"gnl_order_channel"`
	GnlSoCodeApps      string          `json:"gnl_so_code_apps"`
	GnlTotalWeight     float64         `json:"gnl_totalweight"`
	UserID             string          `json:"userid"`
	VoucherApply       []*VoucherApply `json:"voucher_apply"`
	AmountReceived     *AmountReceived `json:"amount_received"`
	DetailItems        []*DetailItem   `json:"detailitems"`
}

type GetSalesInvoiceGPRequest struct {
	Limit    int32  `query:"limit"`
	Offset   int32  `query:"offset"`
	Sopnumbe string `query:"sopnumbe"`
}

type CreateSalesInvoiceRequest struct {
	Code       string `json:"-"`
	CustomerID string `json:"customer_id" valid:"required"`
	// AddressID          string `json:"address_id" valid:"required"`
	// RegionID           string `json:"region_id" valid:"required"`
	DeliveryDateStr string `json:"delivery_date" valid:"required"`
	WrtID           string `json:"wrt_id" valid:"required"`
	// SiteID             string `json:"site_id" valid:"required"`
	RecognitionDateStr string `json:"order_date" valid:"required"`
	// OrderTypeID            string          `json:"order_type_id" valid:"required"`
	// SalespersonID          string          `json:"salesperson_id" valid:"required"`
	SalesTermID string `json:"term_payment_sls_id"`
	// InvoiceTermID          string          `json:"term_invoice_sls_id" valid:"required"`
	// PaymentGroupID         string          `json:"payment_group_id" valid:"required"`
	// ShippingAddress        string          `json:"shipping_address" valid:"required"`
	// BillingAddress         string          `json:"billing_address" valid:"required"`
	RedeemCode             string          `json:"redeem_code"`
	Note                   string          `json:"note"`
	TotalPrice             float64         `json:"-"`
	TotalWeight            float64         `json:"-"`
	TotalCharge            float64         `json:"-"`
	SameTagCustomer        string          `json:"-"`
	DeliveryDate           time.Time       `json:"-"`
	RecognitionDate        time.Time       `json:"-"`
	IsCreateMerchantVa     map[string]int8 `json:"-"`
	CurrentTime            time.Time       `json:"-"`
	CreditLimitBefore      float64         `json:"-"`
	CreditLimitAfter       float64         `json:"-"`
	IsCreateCreditLimitLog bool            `json:"-"`
	NotePriceChange        string          `json:"-"`

	Products []*salesInvoiceItem `json:"products" valid:"required"`
}

type salesInvoiceItem struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"qty"`
	UnitPrice int64   `json:"unit_price"`
	Note      string  `json:"note"`
	// ProductPush   int8    `json:"product_push"`
	TaxPercentage float64 `json:"-"`
	TaxableItem   int8    `json:"-"`
	Subtotal      float64 `json:"-"`
	Weight        float64 `json:"-"`
	DefaultPrice  float64 `json:"-"`
}

// Order Performance

type OrderPerformance struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	QtySell     float64 `json:"qty_sell"`
	AvgSales    float64 `json:"avg_sales"`
	OrderTotal  int64   `json:"order_total"`
}
