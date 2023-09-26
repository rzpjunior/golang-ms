package dto

type PurchaseOrderItemResponse struct {
	ID                 int64   `json:"id"`
	PurchaseOrderID    int64   `json:"purchase_order_id"`
	PurchasePlanItemID int64   `json:"purchase_plan_item_id"`
	ItemID             int64   `json:"item_id"`
	OrderQty           float64 `json:"order_qty"`
	UnitPrice          float64 `json:"unit_price"`
	TaxableItem        int32   `json:"taxable_item"`
	IncludeTax         int32   `json:"include_tax"`
	TaxPercentage      float64 `json:"tax_percentage"`
	TaxAmount          float64 `json:"tax_amount"`
	UnitPriceTax       float64 `json:"unit_price_tax"`
	Subtotal           float64 `json:"subtotal"`
	Weight             float64 `json:"weight"`
	Note               string  `json:"note"`
	PurchaseQty        float64 `json:"purchase_qty"`

	Item *ItemResponse `json:"item"`
}
