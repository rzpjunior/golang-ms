package dto

type SalesInvoiceItemResponse struct {
	ID               int64   `json:"-"`
	SalesInvoiceID   int64   `json:"sales_invoice_id"`
	SalesOrderItemID int64   `json:"sales_order_item_id"`
	ItemID           int64   `json:"item_id"`
	InvoiceQty       float64 `json:"invoice_qty"`
	UnitPrice        float64 `json:"unit_price"`
	Subtotal         float64 `json:"subtotal"`
	Note             string  `json:"note"`
	TaxableItem      int8    `json:"taxable_item"`
	TaxPercentage    float64 `json:"tax_percentage"`
	SkuDiscAmount    float64 `json:"sku_disc_amount"`
}
