package dto

type SalesOrderItemResponse struct {
	ID            int64   `json:"id"`
	SalesOrderId  int64   `json:"sales_order_id"`
	ItemId        int64   `json:"item_id"`
	OrderQty      float64 `json:"order_qty"`
	DefaultPrice  float64 `json:"default_price"`
	UnitPrice     float64 `json:"unit_price"`
	TaxableItem   int32   `json:"taxable_item"`
	TaxPercentage float64 `json:"tax_percentage"`
	ShadowPrice   float64 `json:"shadow_price"`
	Subtotal      float64 `json:"subtotal"`
	Weight        float64 `json:"weight"`
	Note          string  `json:"note"`
}
