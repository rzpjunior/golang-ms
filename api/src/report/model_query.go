package report

type reportDeliveryOrderItem struct {
	OrderCode         string  `orm:"column(order_code)" json:"order_code"`
	DeliveryCode      string  `orm:"column(delivery_code)" json:"delivery_code"`
	ProductCode       string  `orm:"column(product_code)" json:"product_code"`
	ProductName       string  `orm:"column(product_name)" json:"product_name"`
	Uom               string  `orm:"column(uom)" json:"uom"`
	DeliveryItemNote  string  `orm:"column(delivery_item_note)" json:"delivery_item_note"`
	DeliveredQty      float64 `orm:"column(delivered_qty)" json:"delivered_qty"`
	ReceivedQty       float64 `orm:"column(received_qty)" json:"received_qty"`
	DeliveryWeight    float64 `orm:"column(delivery_weight)" json:"delivery_weight"`
	Area              string  `orm:"column(area)" json:"area"`
	Warehouse         string  `orm:"column(warehouse)" json:"warehouse"`
	OrderDeliveryDate string  `orm:"column(order_delivery_date)" json:"order_delivery_date"`
	DeliveryDate      string  `orm:"column(delivery_date)" json:"delivery_date"`
	Wrt               string  `orm:"column(wrt)" json:"wrt"`
	DeliveryStatus    string  `orm:"column(delivery_status)" json:"delivery_status"`
	DeliveryNote      string  `orm:"column(delivery_note)" json:"delivery_note"`
}

type reportPackingOrder struct {
	DeliveryDate   string  `orm:"column(delivery_date)" json:"delivery_date"`
	ProductName    string  `orm:"column(product_name)" json:"product_name"`
	Uom            string  `orm:"column(uom)" json:"uom"`
	TotalOrder     float64 `orm:"column(total_order)" json:"total_order"`
	SubtotalPack   float64 `orm:"column(subtotal_pack)" json:"subtotal_pack"`
	SubtotalWeight float64 `orm:"column(subtotal_weight)" json:"subtotal_weight"`
	HelperCode     string  `orm:"column(helper_code)" json:"helper_code"`
	HelperName     string  `orm:"column(helper_name)" json:"helper_name"`
}

// reportInboundItem: model to contain report Pricing Inbound Item data from query
type reportPricingInboundItem struct {
	InboundCode          string  `orm:"column(inbound_code)" json:"inbound_code"`
	SupplierName         string  `orm:"column(supplier_name)" json:"supplier_name"`
	SupplierCode         string  `orm:"column(supplier_code)" json:"supplier_code"`
	WarehouseOrigin      string  `orm:"column(warehouse_origin)" json:"warehouse_origin"`
	WarehouseDestination string  `orm:"column(warehouse_destination)" json:"warehouse_destination"`
	Area                 string  `orm:"column(area)" json:"area"`
	OrderDate            string  `orm:"column(order_date)" json:"order_date"`
	EtaDate              string  `orm:"column(eta_date)" json:"eta_date"`
	AtaDate              string  `orm:"column(ata_date)" json:"ata_date"`
	ProductCode          string  `orm:"column(product_code)" json:"product_code"`
	ProductName          string  `orm:"column(product_name)" json:"product_name"`
	Uom                  string  `orm:"column(uom)" json:"uom"`
	RequestQty           float64 `orm:"column(request_qty)" json:"request_qty"`
	DeliveredQty         float64 `orm:"column(delivered_qty)" json:"delivered_qty"`
	ReceiveQty           float64 `orm:"column(receive_qty)" json:"receive_qty"`
	InvoiceQty           float64 `orm:"column(invoice_qty)" json:"invoice_qty"`
	Taxability           float64 `orm:"column(taxability)" json:"taxability"`
	TaxPercentage        float64 `orm:"column(tax_percentage)" json:"tax_percentage"`
	UnitPrice            float64 `orm:"column(unit_price)" json:"unit_price"`
	InboundStatus        string  `orm:"column(inbound_status)" json:"inbound_status"`
}

type reportPriceChangeHistory struct {
	CreatedAt   string  `orm:"column(created_at)" json:"created_at"`
	PriceSet    string  `orm:"column(price_set)" json:"price_set"`
	ProductName string  `orm:"column(product_name)" json:"product_name"`
	UnitPrice   float64 `orm:"column(unit_price)" json:"unit_price"`
	CreatedBy   string  `orm:"column(created_by)" json:"created_by"`
}
