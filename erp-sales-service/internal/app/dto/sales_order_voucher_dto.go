package dto

import "time"

type SalesOrderVoucherResponse struct {
	ID           int64     `json:"id"`
	SalesOrderID int64     `json:"sales_order_id"`
	VoucherIDGP  string    `json:"voucher_id_gp"`
	DiscAmount   float64   `json:"disc_amount"`
	CreatedAt    time.Time `json:"created_at"`
	VoucherType  int8      `json:"voucher_type"`
}

type GetSalesOrderVoucherListRequest struct {
	Limit        int64 `json:"limit"`
	Offset       int64 `json:"offset"`
	SalesOrderID int64 `json:"sales_order_id"`
}
