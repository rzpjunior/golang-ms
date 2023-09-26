package dto

import "time"

type PriceTieringLogResponse struct {
	ID               int64     `json:"id"`
	PriceTieringIDGP string    `json:"price_tiering_id_gp"`
	CustomerID       int64     `json:"customer_id"`
	AddressIDGP      string    `json:"address_id_gp"`
	SalesOrderIDGP   string    `json:"sales_order_id_gp"`
	ItemID           int64     `json:"item_id"`
	DiscountQty      float64   `json:"discount_qty"`
	DiscountAmount   float64   `json:"discount_amount"`
	CreatedAt        time.Time `json:"created_at"`
	Status           int8      `json:"status"`
}

type PriceTieringLogRequestCreate struct {
	PriceTieringIDGP string  `json:"price_tiering_id_gp"`
	CustomerID       int64   `json:"customer_id"`
	AddressIDGP      string  `json:"address_id_gp"`
	SalesOrderIDGP   string  `json:"sales_order_id_gp"`
	ItemID           int64   `json:"item_id"`
	DiscountQty      float64 `json:"discount_qty"`
	DiscountAmount   float64 `json:"discount_amount"`
}

type PriceTieringLogRequestGet struct {
	Search           string `json:"search"`
	PriceTieringIDGP string `json:"price_tiering_id_gp"`
	CustomerID       int64  `json:"customer_id"`
	AddressIDGP      string `json:"address_id_gp"`
	SalesOrderIDGP   string `json:"sales_order_id_gp"`
	ItemID           int64  `json:"item_id"`
	OrderBy          string `json:"order_by"`
	Offset           int64  `json:"offset"`
	Limit            int64  `json:"limit"`
	Status           int8   `json:"status"`
}

type PriceTieringLogRequestCancel struct {
	PriceTieringIDGP string `json:"price_tiering_id_gp"`
	CustomerID       int64  `json:"customer_id"`
	AddressIDGP      string `json:"address_id_gp"`
	SalesOrderIDGP   string `json:"sales_order_id_gp"`
	ItemID           int64  `json:"item_id"`
}
