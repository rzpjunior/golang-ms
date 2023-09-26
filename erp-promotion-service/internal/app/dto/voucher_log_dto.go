package dto

import "time"

type VoucherLogResponse struct {
	ID                    int64     `json:"id"`
	VoucherID             int64     `json:"voucher_id"`
	CustomerID            int64     `json:"customer_id"`
	AddressIDGP           string    `json:"address_id_gp"`
	SalesOrderIDGP        string    `json:"sales_order_id_gp"`
	VoucherDiscountAmount float64   `json:"voucher_discount_amount"`
	Status                int8      `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
}

type VoucherLogRequestCreate struct {
	VoucherID             int64   `json:"voucher_id"`
	CustomerID            int64   `json:"customer_id"`
	AddressIDGP           string  `json:"address_id_gp"`
	SalesOrderIDGP        string  `json:"sales_order_id_gp"`
	VoucherDiscountAmount float64 `json:"voucher_discount_amount"`
}

type VoucherLogRequestGet struct {
	Search         string `json:"search"`
	CustomerID     int64  `json:"customer_id"`
	SalesOrderIDGP string `json:"sales_order_id_gp"`
	AddressIDGP    string `json:"address_id_gp"`
	Status         int8   `json:"status"`
	VoucherID      int64  `json:"voucher_id"`
	OrderBy        string `json:"order_by"`
	Offset         int64  `json:"offset"`
	Limit          int64  `json:"limit"`
	CreatedDate    string `json:"created_date"`
	Code           string `json:"code"`
}

type VoucherLogRequestCancel struct {
	ID             int64  `json:"id"`
	VoucherID      int64  `json:"voucher_id"`
	CustomerID     int64  `json:"customer_id"`
	AddressIDGP    string `json:"address_id_gp"`
	SalesOrderIDGP string `json:"sales_order_id_gp"`
	Status         int8   `json:"status"`
	Code           string `json:"code"`
}
