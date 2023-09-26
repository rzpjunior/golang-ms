package dto

type DeliveryKoliResponse struct {
	SalesOrderCode string  `json:"sales_order_code"`
	KoliId         int64   `json:"koli_id"`
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
}
