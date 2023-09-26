package dto

type SalesInvoiceItemResponse struct {
	ID               int64   `json:"-"`
	SalesInvoiceID   string  `json:"sales_invoice_id"`
	SalesOrderItemID int64   `json:"sales_order_item_id"`
	ItemID           string  `json:"item_id"`
	ItemName         string  `json:"name"`
	InvoiceQty       float64 `json:"invoice_qty"`
	UnitPrice        float64 `json:"unit_price"`
	Subtotal         float64 `json:"subtotal"`
	Note             string  `json:"note"`
	TaxableItem      int8    `json:"taxable_item"`
	TaxPercentage    float64 `json:"tax_percentage"`
	SkuDiscAmount    float64 `json:"sku_disc_amount"`
	UomName          string  `json:"uom_name"`

	Item *ItemResponse `json:"item"`
}

type SalesInvoiceItemListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type SalesInvoiceItemDetailRequest struct {
	Id int32 `json:"id"`
}
