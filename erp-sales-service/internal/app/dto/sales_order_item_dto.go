package dto

type SalesOrderItemResponse struct {
	ID               int64   `json:"-"`
	SalesOrderID     int64   `json:"sales_order_id"`
	ItemIDGP         string  `json:"item_id_gp"`
	ItemName         string  `json:"item_name"`
	PriceTieringIDGP string  `json:"price_tiering_id_gp"`
	OrderQty         float64 `json:"order_qty"`
	UnitPrice        float64 `json:"unit_price"`
	UomIDGP          string  `json:"uom_gp"`
	UomName          string  `json:"uom_name"`
	Subtotal         float64 `json:"subtotal"`
	Weight           float64 `json:"weight"`
	ImageUrl         string  `json:"image_url"`
}
