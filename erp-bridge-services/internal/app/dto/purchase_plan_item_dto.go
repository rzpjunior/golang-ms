package dto

type PurchasePlanItemResponse struct {
	ID              int64   `json:"id"`
	PurchasePlanID  int64   `json:"purchase_plan_id"`
	ItemID          int64   `json:"item_id"`
	PurchasePlanQty float64 `json:"purchase_plan_qty"`
	PurchaseQty     float64 `json:"purchase_qty"`
	UnitPrice       float64 `json:"unit_price"`
	Subtotal        float64 `json:"subtotal"`
	Weight          float64 `json:"weight"`
}
