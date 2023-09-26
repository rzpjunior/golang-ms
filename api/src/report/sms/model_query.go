// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import "time"

// reportPurchaseOrder : model to contain report purchase order data from query
type reportPurchaseOrder struct {
	SupplierCode  string  `orm:"column(supplier_code)" json:"supplier_code"`
	SupplierName  string  `orm:"column(supplier_name)" json:"supplier_name"`
	WarehouseName string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	OrderCode     string  `orm:"column(order_code)" json:"order_code"`
	OrderDate     string  `orm:"column(order_date)" json:"order_date"`
	EtaDate       string  `orm:"column(eta_date)" json:"eta_date"`
	EtaTime       string  `orm:"column(eta_time)" json:"eta_time"`
	OrderStatus   string  `orm:"column(order_status)" json:"order_status"`
	PaymentTerm   string  `orm:"column(order_payment_term)" json:"order_payment_term"`
	DeliveryFee   float64 `orm:"column(delivery_fee)" json:"delivery_fee"`
	TaxAmount     float64 `orm:"column(tax_amount)" json:"tax_amount"`
	GrandTotal    float64 `orm:"column(grand_total)" json:"grand_total"`
	OrderNote     string  `orm:"column(order_note)" json:"order_note"`
	TotalPrice    float64 `orm:"column(total_price)" json:"total_price"`
	SupplierBadge string  `orm:"column(supplier_badge)" json:"supplier_badge"`
	GoodsReceipt  string  `orm:"column(good_receipt)" json:"goods_receipt"`
}

// reportPurchaseOrderItem : model to contain report purchase order item data from query
type reportPurchaseOrderItem struct {
	OrderCode          string  `orm:"column(order_code)" json:"order_code"`
	ProductCode        string  `orm:"column(product_code)" json:"product_code"`
	ProductName        string  `orm:"column(product_name)" json:"product_name"`
	UOM                string  `orm:"column(uom)" json:"uom"`
	OrderItemNote      string  `orm:"column(order_item_note)" json:"order_item_note"`
	OrderedQty         float64 `orm:"column(ordered_qty)" json:"ordered_qty"`
	OrderUnitPrice     float64 `orm:"column(order_unit_price)" json:"order_unit_price"`
	TaxableStr         string  `orm:"column(taxable_item_str)" json:"taxable_item_str"`
	IncludeTaxStr      string  `orm:"column(include_tax_str)" json:"include_tax_str"`
	OrderTaxPercentage float64 `orm:"column(order_tax_percentage)" json:"order_tax_percentage"`
	OrderTaxAmount     float64 `orm:"column(order_tax_amount)" json:"order_tax_amount"`
	OrderUnitPriceTax  float64 `orm:"column(order_unit_price_tax)" json:"order_unit_price_tax"`
	Subtotal           float64 `orm:"column(subtotal)" json:"subtotal"`
	TotalWeight        float64 `orm:"column(total_weight)" json:"total_weight"`
	AreaName           string  `orm:"column(area_name)" json:"area_name"`
	WarehouseName      string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	SupplierCode       string  `orm:"column(supplier_code)" json:"supplier_code"`
	SupplierName       string  `orm:"column(supplier_name)" json:"supplier_name"`
	OrderDate          string  `orm:"column(order_date)" json:"order_date"`
	EtaDate            string  `orm:"column(eta_date)" json:"eta_date"`
	InvoicedQty        float64 `orm:"column(invoiced_qty)" json:"invoiced_qty"`
	PurchaseQty        float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
}

// reportPurchaseInvoice : model to contain report purchase invoice data from query
type reportPurchaseInvoice struct {
	PurchaseInvoiceID int64     `orm:"column(invoice_id)" json:"-"`
	AreaName          string    `orm:"column(area_name)" json:"area_name"`
	WarehouseName     string    `orm:"column(warehouse_name)" json:"warehouse_name"`
	OrderCode         string    `orm:"column(order_code)" json:"order_code"`
	OrderDate         string    `orm:"column(order_date)" json:"order_date"`
	EtaDate           string    `orm:"column(eta_date)" json:"eta_date"`
	SupplierCode      string    `orm:"column(supplier_code)" json:"supplier_code"`
	SupplierName      string    `orm:"column(supplier_name)" json:"supplier_name"`
	TotalOrder        float64   `orm:"column(total_order)" json:"total_order"`
	InvoiceCode       string    `orm:"column(invoice_code)" json:"invoice_code"`
	InvoiceDate       string    `orm:"column(invoice_date)" json:"invoice_date"`
	InvoiceDueDate    string    `orm:"column(invoice_due_date)" json:"invoice_due_date"`
	InvoiceStatus     string    `orm:"column(invoice_status)" json:"invoice_status"`
	InvoiceNote       string    `orm:"column(invoice_note)" json:"invoice_note"`
	DeliveryFee       float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	InvoiceAmount     float64   `orm:"column(invoice_amount)" json:"invoice_amount"`
	TaxAmount         float64   `orm:"column(tax_amount)" json:"tax_amount"`
	TotalInvoice      float64   `orm:"column(total_invoice)" json:"total_invoice"`
	AdjustmentAmount  float64   `orm:"column(adjustment_amount)" json:"adjustment_amount"`
	AdjustmentNote    string    `orm:"column(adjustment_note)" json:"adjustment_note"`
	CreatedAt         string    `orm:"column(created_at)" json:"created_at"`
	CreatedBy         string    `orm:"column(created_by)" json:"created_by"`
	SupplierType      string    `orm:"column(supplier_type)" json:"supplier_type"`
	PaymentTerm       string    `orm:"column(payment_term)" json:"payment_term"`
	TotalPayment      float64   `orm:"column(total_payment)" json:"total_payment"`
	AtaDate           time.Time `orm:"column(ata_date)" json:"ata_date"`
	LastUpdatedAt     time.Time `orm:"column(last_updated_at)" json:"last_updated_at"`
	LastUpdatedBy     string    `orm:"column(last_updated_by)" json:"last_updated_by"`
}

// reportPurchasePayment : model to contain report purchase payment data from query
type reportPurchasePayment struct {
	AreaName      string  `orm:"column(area_name)" json:"area_name"`
	SupplierCode  string  `orm:"column(supplier_code)" json:"supplier_code"`
	SupplierName  string  `orm:"column(supplier_name)" json:"supplier_name"`
	PaymentCode   string  `orm:"column(payment_code)" json:"payment_code"`
	PaymentDate   string  `orm:"column(payment_date)" json:"payment_date"`
	PaymentAmount float64 `orm:"column(payment_amount)" json:"payment_amount"`
	PaymentMethod string  `orm:"column(payment_method)" json:"payment_method"`
	PaymentStatus string  `orm:"column(payment_status)" json:"payment_status"`
	InvoiceCode   string  `orm:"column(invoice_code)" json:"invoice_code"`
	TotalInvoice  float64 `orm:"column(total_invoice)" json:"total_invoice"`
	PaymentNumber string  `orm:"column(payment_number)" json:"bank_payment_voucher_number"`
	CreatedAt     string  `orm:"column(created_at)" json:"created_at"`
	CreatedBy     string  `orm:"column(created_by)" json:"created_by"`
}

// reportPurchaseInvoiceItem : model to contain report purchase invoice item data from query
type reportPurchaseInvoiceItem struct {
	SupplierName  string  `orm:"column(supplier_name)" json:"supplier_name"`
	WarehouseName string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	Area          string  `orm:"column(area)" json:"area"`
	OrderCode     string  `orm:"column(order_code)" json:"order_code"`
	OrderStatus   string  `orm:"column(order_status)" json:"order_status"`
	InvoiceCode   string  `orm:"column(invoice_code)" json:"invoice_code"`
	InvoiceStatus string  `orm:"column(invoice_status)" json:"invoice_status"`
	GRCode        string  `orm:"column(gr_code)" json:"gr_code"`
	GRStatus      string  `orm:"column(gr_status)" json:"gr_status"`
	EtaDate       string  `orm:"column(eta_date)" json:"eta_date"`
	ProductCode   string  `orm:"column(product_code)" json:"product_code"`
	ProductName   string  `orm:"column(product_name)" json:"product_name"`
	UOM           string  `orm:"column(uom)" json:"uom"`
	UnitPrice     float64 `orm:"column(unit_price)" json:"unit_price"`
	OrderQty      float64 `orm:"column(order_qty)" json:"order_qty"`
	TaxableStr    string  `orm:"column(taxable_item_str)" json:"taxable_item_str"`
	IncludeTaxStr string  `orm:"column(include_tax_str)" json:"include_tax_str"`
	TaxPercentage float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	TaxAmount     float64 `orm:"column(tax_amount)" json:"tax_amount"`
	UnitPriceTax  float64 `orm:"column(unit_price_tax)" json:"unit_price_tax"`

	DeliveredQty float64 `orm:"column(delivered_qty)" json:"delivered_qty"`
	ReceivedQty  float64 `orm:"column(received_qty)" json:"received_qty"`
	InvoiceQty   float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	RejectQty    float64 `orm:"column(reject_qty)" json:"reject_qty"`
	DeliveryFee  float64 `orm:"column(delivery_fee)" json:"delivery_fee"`
	TotalInvoice float64 `orm:"column(total_invoice)" json:"total_invoice"`
}

// reportCogs : model to contain report cogs data from query
type reportCogs struct {
	AreaName      string  `orm:"column(area_name)" json:"area_name"`
	WarehouseName string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	EtaDate       string  `orm:"column(eta_date)" json:"eta_date"`
	ProductCode   string  `orm:"column(product_code)" json:"product_code"`
	ProductName   string  `orm:"column(product_name)" json:"product_name"`
	UOM           string  `orm:"column(uom)" json:"uom"`
	AvgPrice      float64 `orm:"column(avg_price)" json:"avg_price"`
}

// reportPriceComparison : model to contain report price comparison data from query
type reportPriceComparison struct {
	SurveyDate   string  `orm:"column(survey_date)"`
	AreaName     string  `orm:"column(area_name)"`
	ProductCode  string  `orm:"column(product_code)"`
	ProductName  string  `orm:"column(product_name)"`
	UOM          string  `orm:"column(uom)"`
	SellingPrice float64 `orm:"column(selling_price)"`
	PublicPrice1 float64 `orm:"column(public_price_1)"`
	PublicPrice2 float64 `orm:"column(public_price_2)"`
}

// reportInboundDetail : model to contain report inbound detail data from query
type reportInboundDetail struct {
	SupplierName  string    `orm:"column(supplier_name)" json:"supplier_name"`
	WarehouseName string    `orm:"column(warehouse_name)" json:"warehouse_name"`
	OrderCode     string    `orm:"column(po_code)" json:"po_code"`
	EtaDate       time.Time `orm:"column(eta_date);type(date)" json:"eta_date"`
	EtaTime       string    `orm:"column(eta_time)" json:"eta_time"`
	CommittedAt   time.Time `orm:"column(committed_at);type(datetime)" json:"committed_at"`
	AtaDate       time.Time `orm:"column(ata_date);type(date)" json:"ata_date"`
	AtaTime       string    `orm:"column(ata_time)" json:"ata_time"`
	Source        string    `orm:"column(source)" json:"-"`
}

// reportInboundSummary : model to contain report inbound summary data
type reportInboundSummary struct {
	Source              string
	EtaDate             time.Time
	TotalData           int64
	TotalFulfillInbound int64
	TotalFulfillCommit  int64
}

// reportFieldPurchaser : model to contain report field purchaser data from query
type reportFieldPurchaser struct {
	PurchasePlanDate         string  `orm:"column(purchase_plan_date)" json:"purchase_plan_date"`
	PurchasePlanCode         string  `orm:"column(purchase_plan_code)" json:"purchase_plan_code"`
	PurchaseOrderDate        string  `orm:"column(purchase_order_date)" json:"purchase_order_date"`
	PurchaseOrderCode        string  `orm:"column(purchase_order_code)" json:"purchase_order_code"`
	SupplierOrganizationName string  `orm:"column(supplier_organization_name)" json:"supplier_organization_name"`
	SupplierName             string  `orm:"column(supplier_name)" json:"supplier_name"`
	Returnable               string  `orm:"column(returnable)" json:"returnable"`
	Rejectable               string  `orm:"column(rejectable)" json:"rejectable"`
	ProductCode              string  `orm:"column(product_code)" json:"product_code"`
	ProductName              string  `orm:"column(product_name)" json:"product_name"`
	UOM                      string  `orm:"column(uom)" json:"uom"`
	PurchasePlanQty          float64 `orm:"column(purchase_plan_qty" json:"purchase_plan_qty"`
	PriceReference           float64 `orm:"column(price_reference)" json:"price_reference"`
	PurchaseQty              float64 `orm:"column(purchase_qty)" json:"purchase_qty"`
	UnitPrice                float64 `orm:"column(unit_price)" json:"unit_price"`
	TotalPrice               float64 `orm:"column(total_price)" json:"total_price"`
	PaymentTerm              string  `orm:"column(payment_term)" json:"payment_term"`
	FieldPurchaserName       string  `orm:"column(field_purchaser_name)" json:"field_purchaser_name"`
	OrderLocation            string  `orm:"column(order_location)" json:"order_location"`
	ConsolidatedShipmentCode string  `orm:"column(consolidated_shipment_code)" json:"consolidated_shipment_code"`
	WarehouseName            string  `orm:"column(warehouse_name)" json:"warehouse_name"`
	DriverName               string  `orm:"column(driver_name)" json:"driver_name"`
	VehicleNumber            string  `orm:"column(vehicle_number)" json:"vehicle_number"`
	DriverPhoneNumber        string  `orm:"column(driver_phone_number)" json:"driver_phone_number"`
	EtaDate                  string  `orm:"column(eta_date)" json:"eta_date"`
	EtaTime                  string  `orm:"column(eta_time)" json:"eta_time"`
	AtaDate                  string  `orm:"column(ata_date)" json:"ata_date"`
	AtaTime                  string  `orm:"column(ata_time)" json:"ata_time"`
	InboundDate              string  `orm:"column(inbound_date)" json:"inbound_date"`
	ReceiveQty               float64 `orm:"column(receive_qty)" json:"receive_qty"`
	Status                   string  `orm:"column(status)" json:"status"`
}
