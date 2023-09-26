package report

import "time"

type reportSalesOrder struct {
	SalesOrderCode         string  `orm:"column(sales_order_code)" json:"sales_order_code"`
	CustomerCode           string  `orm:"column(customer_code)" json:"customer_code"`
	CustomerTag            string  `orm:"column(customer_tag)" json:"customer_tag"`
	CustomerName           string  `orm:"column(customer_name)" json:"customer_name"`
	CustomerPhoneNumber    string  `orm:"column(customer_phone_number)" json:"customer_phone_number"`
	RecipientName          string  `orm:"column(recipient_name)" json:"recipient_name"`
	RecipientPhoneNumber   string  `orm:"column(recipient_phone_number)" json:"recipient_phone_number"`
	ShippingAddress        string  `orm:"column(shipping_address)" json:"shipping_address"`
	WarehouseName          string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	ArchetypeName          string  `orm:"column(archetype_name)" json:"archetype_name"`
	AreaName               string  `orm:"column(area_name)" json:"area_name"`
	Salesperson            string  `orm:"column(salesperson)" json:"salesperson"`
	SalesGroup             string  `orm:"column(sales_group)" json:"sales_group"`
	OrderDate              string  `orm:"column(order_date)" json:"order_date"`
	OrderDeliveryDate      string  `orm:"column(order_delivery_date)" json:"order_delivery_date"`
	OrderStatus            string  `orm:"column(order_status)" json:"order_status"`
	OrderNote              string  `orm:"column(order_note)" json:"order_note"`
	TotalSKUDiscountAmount float64 `orm:"column(total_sku_disc_amount)" json:"total_sku_disc_amount"`
	GrandTotal             float64 `orm:"column(grand_total)" json:"grand_total"`
	OrderChannel           string  `orm:"column(order_channel)" json:"order_channel"`
	PromoCode              string  `orm:"column(promo_code)" json:"promo_code"`
	DeliveryFee            float64 `orm:"column(delivery_fee)" json:"delivery_fee"`
	OrderType              string  `orm:"column(order_type_name)" json:"order_type"`
	CancelType             string  `orm:"column(cancel_type)" json:"cancel_type"`
	ETD                    int64   `orm:"column(estimate_time_departure)"`

	BusinessType  string `orm:"column(business_type)" json:"business_type"`
	City          string `orm:"column(city)" json:"city"`
	CreatedAt     string `orm:"column(created_at)" json:"created_at"`
	CreatedBy     string `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt string `orm:"column(last_updated_at)" json:"last_updated_at"`
	LastUpdatedBy string `orm:"column(last_updated_by)" json:"last_updated_by"`

	ETDSt string `json:"estimate_time_departure"`
}

type reportSalesOrderItem struct {
	SalesOrderCode       string  `orm:"column(sales_order_code)" json:"sales_order_code"`
	ProductCode          string  `orm:"column(product_code)" json:"product_code"`
	ProductName          string  `orm:"column(product_name)" json:"product_name"`
	CategoryName         string  `orm:"column(category_name)" json:"category_name"`
	UomName              string  `orm:"column(uom_name)" json:"uom_name"`
	OrderItemNote        string  `orm:"column(order_item_note)" json:"order_item_note"`
	OrderedQty           float64 `orm:"column(ordered_qty)" json:"ordered_qty"`
	InvoiceQty           float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	OrderUnitPrice       float64 `orm:"column(order_unit_price)" json:"order_unit_price"`
	OrderUnitShadowPrice float64 `orm:"column(order_unit_shadow_price)" json:"order_unit_shadow_price"`
	TaxableStr           string  `orm:"column(taxable_item_str)" json:"taxable_item_str"`
	OrderTaxPercentage   float64 `orm:"column(order_tax_percentage)" json:"order_tax_percentage"`
	Subtotal             float64 `orm:"column(subtotal)" json:"subtotal"`
	TotalWeight          float64 `orm:"column(total_weight)" json:"total_weight"`
	OrderDate            string  `orm:"column(order_date)" json:"order_date"`
	OrderDeliveryDate    string  `orm:"column(order_delivery_date)" json:"order_delivery_date"`
	AreaName             string  `orm:"column(area_name)" json:"area_name"`
	WarehouseName        string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	WrtName              string  `orm:"column(wrt_name)" json:"wrt_name"`
	OrderStatus          string  `orm:"column(order_status)" json:"order_status"`
	DiscountQty          float64 `orm:"column(discount_qty)" json:"discount_qty"`
	UnitPriceDiscount    float64 `orm:"column(unit_price_discount)" json:"unit_price_discount"`
	SkuDiscAmount        float64 `orm:"column(sku_disc_amount)" json:"sku_disc_amount"`
	SKUDiscountName      string  `orm:"column(sku_discount_name)" json:"sku_discount_name"`
	SalesOrderTypeName   string  `orm:"column(sales_order_type_name)" json:"sales_order_type_name"`
}

type reportSalesInvoice struct {
	SalesOrderCode     string  `orm:"column(order_code)" json:"order_code"`
	MerchantCode       string  `orm:"column(merchant_code)" json:"merchant_code"`
	MerchantName       string  `orm:"column(merchant_name)" json:"merchant_name"`
	BranchCode         string  `orm:"column(branch_code)" json:"branch_code"`
	BranchName         string  `orm:"column(branch_name)" json:"branch_name"`
	OrderDeliveryDate  string  `orm:"column(order_delivery_date)" json:"order_delivery_date"`
	InvoiceCode        string  `orm:"column(invoice_code)" json:"invoice_code"`
	InvoiceDate        string  `orm:"column(invoice_date)" json:"invoice_date"`
	InvoiceDueDate     string  `orm:"column(invoice_due_date)" json:"invoice_due_date"`
	InvoiceStatus      string  `orm:"column(invoice_status)" json:"invoice_status"`
	AdjNote            string  `orm:"column(adjustment_note)" json:"adjustment_note"`
	TotalConfPay       float64 `orm:"column(total_confirmed_payment)" json:"total_confirmed_payment"`
	TotalInvoice       float64 `orm:"column(total_invoice)" json:"total_invoice"`
	DeliveryFee        float64 `orm:"column(delivery_fee)" json:"delivery_fee"`
	VouAmount          float64 `orm:"column(voucher_amount)" json:"voucher_amount"`
	AdjAmount          float64 `orm:"column(adjustment_amount)" json:"adjustment_amount"`
	TotalCharge        float64 `orm:"column(total_charge)" json:"total_charge"`
	Area               string  `orm:"column(area)" json:"area"`
	Warehouse          string  `orm:"column(warehouse)" json:"warehouse"`
	CustomerGroup      string  `orm:"column(customer_group)" json:"customer_group"`
	BusinessType       string  `orm:"column(business_type)" json:"business_type"`
	Archetype          string  `orm:"column(archetype)" json:"archetype"`
	PaymentTerm        string  `orm:"column(payment_term)" json:"payment_term"`
	InvoiceTerm        string  `orm:"column(invoice_term)" json:"invoice_term"`
	PointRedeemAmount  float64 `orm:"column(point_redeem_amount)" json:"point_redeem_amount"`
	CreatedAt          string  `orm:"column(created_at)" json:"created_at"`
	CreatedBy          string  `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt      string  `orm:"column(updated_at)" json:"updated_at"`
	LastUpdatedBy      string  `orm:"column(updated_by)" json:"updated_by"`
	TotalSkuDiscAmount float64 `orm:"column(total_sku_disc_amount)" json:"total_sku_disc_amount"`
}

type reportSalesPayment struct {
	PaymentCode       string  `orm:"column(payment_code)" json:"payment_code"`
	BankReceiveVouNum string  `orm:"column(bank_receive_num)" json:"bank_receive_num"`
	PaymentDate       string  `orm:"column(payment_date)" json:"payment_date"`
	ReceivedDate      string  `orm:"column(received_date);type(date);null" json:"received_date"`
	Area              string  `orm:"column(area)" json:"area"`
	WarehouseName     string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	PaymentStatus     string  `orm:"column(payment_status)" json:"payment_status"`
	PaymentMethod     string  `orm:"column(payment_method)" json:"payment_method"`
	PaymentAmount     float64 `orm:"column(payment_amount)" json:"payment_amount"`
	OrderCode         string  `orm:"column(order_code)" json:"order_code"`
	InvoiceCode       string  `orm:"column(invoice_code)" json:"invoice_code"`
	OutletCode        string  `orm:"column(outlet_code)" json:"outlet_code"`
	OutletName        string  `orm:"column(outlet_name)" json:"outlet_name"`
	CreatedAt         string  `orm:"column(created_at)" json:"created_at"`
	CreatedBy         string  `orm:"column(created_by)" json:"created_by"`
}

type reportProspectiveCustomer struct {
	CustomerCode       string `orm:"column(prospect_customer_code)" json:"prospect_customer_code"`
	CustomerName       string `orm:"column(prospect_customer_name)" json:"prospect_customer_name"`
	BusinessType       string `orm:"column(business_type)" json:"business_type"`
	Archetype          string `orm:"column(archetype)" json:"archetype"`
	PicName            string `orm:"column(pic_name)" json:"pic_name"`
	PhoneNumber        string `orm:"column(phone_number)" json:"phone_number"`
	PicFinanceName     string `orm:"column(pic_finance_name);size(100)" json:"pic_finance_name"`
	PicFinanceContact  string `orm:"column(pic_finance_contact);size(15)" json:"pic_finance_contact"`
	PicBusinessName    string `orm:"column(pic_business_name);size(100);null" json:"pic_business_name"`
	PicBusinessContact string `orm:"column(pic_business_contact);size(15);null" json:"pic_business_contact"`
	InvoiceTerm        string `orm:"column(term_of_invoice);null" json:"invoice_term,omitempty"`
	PaymentTerm        string `orm:"column(term_of_payment);null" json:"payment_term,omitempty"`
	BillingAddress     string `orm:"column(billing_address);size(350);null" json:"billing_address"`
	Note               string `orm:"column(notes);size(250);null" json:"note"`
	Area               string `orm:"column(area)" json:"area"`
	StreetAddress      string `orm:"column(street_address)" json:"street_address"`
	Province           string `orm:"column(province)" json:"province"`
	City               string `orm:"column(city)" json:"city"`
	District           string `orm:"column(district)" json:"district"`
	SubDistrict        string `orm:"column(sub_district)" json:"sub_district"`
	PostalCode         string `orm:"column(postal_code)" json:"postal_code"`
	BestTimeToCall     string `orm:"column(best_time_to_call)" json:"best_time_to_call"`
	ReferralCode       string `orm:"column(referral_code)" json:"referral_code"`
	ExistingCustomer   string `orm:"column(existing_customer)" json:"existing_customer"`
	ReqUpgrade         string `orm:"column(request_upgrade)" json:"request_upgrade"`
	CreatedAt          string `orm:"column(created_at)" json:"created_at"`
	ProcessedAt        string `orm:"column(processed_at)" json:"processed_at"`
	ProcessedBy        string `orm:"column(processed_by)" json:"processed_by"`
	Status             string `orm:"column(status)" json:"status"`
	Salesperson        string `orm:"column(salesperson)" json:"salesperson"`
	DeclineType        string `orm:"column(decline_type)" json:"decline_type"`
	DeclineNote        string `orm:"column(decline_note)" json:"decline_note"`
}

// reportSkuDiscount : struct to hold sku discount data
type reportSkuDiscount struct {
	SkuDiscName  string    `orm:"column(sku_disc_name)" json:"sku_disc"`
	PriceSet     string    `orm:"column(price_set_name)" json:"price_set_name"`
	StartPeriod  time.Time `orm:"column(start_timestamp)" json:"start_period"`
	EndPeriod    time.Time `orm:"column(end_timestamp)" json:"end_period"`
	Division     string    `orm:"column(division_name)" json:"division"`
	OrderChannel string    `orm:"column(order_channel_name)" json:"order_channel"`
	Note         string    `orm:"column(note)" json:"note"`
}

// reportSkuDiscountItem : struct to hold sku discount item data
type reportSkuDiscountItem struct {
	SkuDiscName         string  `orm:"column(sku_disc_name)" json:"sku_disc"`
	ProductCode         string  `orm:"column(product_code)" json:"product_code"`
	ProductName         string  `orm:"column(product_name)" json:"product"`
	Uom                 string  `orm:"column(uom_name)" json:"uom"`
	TierLevel           int8    `orm:"column(tier_level)" json:"tier_level"`
	MinimumQty          float64 `orm:"column(minimum_qty)" json:"minimum_qty"`
	Amount              float64 `orm:"column(disc_amount)" json:"disc_amount"`
	OverallQuota        float64 `orm:"column(overall_quota)" json:"overall_quota"`
	OverallQuotaPerUser float64 `orm:"column(overall_quota_per_user)" json:"overall_quota_per_user"`
	DailyQuotaPerUser   float64 `orm:"column(daily_quota_per_user)" json:"daily_quota_per_user"`
	Budget              float64 `orm:"column(budget)" json:"budget"`
	RemBudget           float64 `orm:"column(rem_budget)" json:"rem_budget"`
}

type reportSalesOrderFeedback struct {
	MerchantCode          string `orm:"column(Merchant_Code)" json:"Merchant_Code"`
	MerchantName          string `orm:"column(Merchant_Name)" json:"Merchant_Name"`
	MerchantPhoneNumber   string `orm:"column(Merchant_Phone_Number)" json:"Merchant_Phone_Number"`
	BusinessType          string `orm:"column(Business_Type)" json:"Business_Type"`
	Archetype             string `orm:"column(Archetype)" json:"Archetype"`
	BranchShippingAddress string `orm:"column(Branch_Shipping_Address)" json:"Branch_Shipping_Address"`
	City                  string `orm:"column(City)" json:"City"`
	District              string `orm:"column(District)" json:"District"`
	Area                  string `orm:"column(Area)" json:"Area"`
	SalesOrderCode        string `orm:"column(Sales_Order_Code)" json:"Sales_Order_Code"`
	DeliveryDate          string `orm:"column(Delivery_Date)" json:"Delivery_Date"`
	FeedBackCreatedAt     string `orm:"column(Feedback_Created_At)" json:"Feedback_Created_At"`
	RatingScore           string `orm:"column(Rating_Score)" json:"Rating_Score"`
	Tags                  string `orm:"column(Tags)" json:"Tags"`
	FeedbackDescription   string `orm:"column(Feedback_Description)" json:"Feedback_Description"`
	ToBeContacted         string `orm:"column(To_Be_Contacted)" json:"To_Be_Contacted"`
}

type reportEdenPointLog struct {
	EdenPointDate      string  `orm:"column(edenpoint_date)" json:"edenpoint_date"`
	MerchantName       string  `orm:"column(merchant_name)" json:"merchant_name"`
	PreviousEdenPoint  float64 `orm:"column(previous_edenpoint)" json:"previous_edenpoint"`
	EdenPoint          float64 `orm:"column(edenpoint)" json:"edenpoint"`
	Status             string  `orm:"column(status)" json:"status"`
	CurrentEdenPoint   float64 `orm:"column(current_edenpoint)" json:"current_edenpoint"`
	TransactionType    string  `orm:"column(transaction_type)" json:"transaction_type"`
	AdvocateMerchant   string  `orm:"column(advocate_merchant)" json:"advocate_merchant"`
	RefereeMerchant    string  `orm:"column(referee_merchant)" json:"referee_merchant"`
	CampaignId         string  `orm:"column(campaign_id)" json:"campaign_id"`
	CampaignName       string  `orm:"column(campaign_name)" json:"campaign_name"`
	CampaignMultiplier int8    `orm:"column(campaign_multiplier)" json:"campaign_multiplier"`
	LogNote            string  `orm:"column(log_note)" json:"log_note"`
	OrderCode          string  `orm:"column(order_code)" json:"order_code"`
	OrderDate          string  `orm:"column(order_date)" json:"order_date"`
	CreatedOrderDate   string  `orm:"column(created_at)" json:"created_at"`
	FinishedOrderDate  string  `orm:"column(finished_at)" json:"finished_at"`
	TotalSalesOrder    float64 `orm:"column(total_sales_order)" json:"total_sales_order"`
	OrderStatus        string  `orm:"column(order_status)" json:"order_status"`
}
