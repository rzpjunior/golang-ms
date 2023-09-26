package dto

type PurchasePlanItemResponse struct {
	ID              string              `json:"id"`
	PurchasePlanID  string              `json:"purchase_plan_id"`
	Item            *ItemResponse       `json:"item"`
	OrderQty        float64             `json:"order_qty"`
	PurchaseQty     float64             `json:"purchase_qty"`
	UnitPrice       float64             `json:"unit_price"`
	Subtotal        float64             `json:"subtotal"`
	Weight          float64             `json:"weight"`
	PurchaseOrderID []PurchaseOrderItem `json:"purchase_order_items"`
}
