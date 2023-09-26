package dto

type PurchaseOrderItemResponse struct {
	ID                 string  `json:"id,omitempty"`
	PurchaseOrderID    string  `json:"purchase_order_id,omitempty"`
	PurchasePlanItemID int64   `json:"purchase_plan_item_id,omitempty"`
	ItemID             string  `json:"item_id,omitempty"`
	ItemName           string  `json:"item_name,omitempty"`
	Uom                string  `json:"uom,omitempty"`
	OrderQty           float64 `json:"order_qty,omitempty"`
	UnitPrice          float64 `json:"unit_price,omitempty"`
	TaxableItem        int32   `json:"taxable_item,omitempty"`
	IncludeTax         int32   `json:"include_tax,omitempty"`
	TaxPercentage      float64 `json:"tax_percentage,omitempty"`
	TaxAmount          float64 `json:"tax_amount,omitempty"`
	UnitPriceTax       float64 `json:"unit_price_tax,omitempty"`
	Subtotal           float64 `json:"subtotal,omitempty"`
	Weight             float64 `json:"weight,omitempty"`
	Note               string  `json:"note"`
	PurchaseQty        float64 `json:"purchase_qty,omitempty"`

	Item *ItemResponse `json:"item,omitempty"`
}

type PurchaseOrderItemListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type PurchaseOrderItemDetailRequest struct {
	Id int32 `json:"id"`
}
