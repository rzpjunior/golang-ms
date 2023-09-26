package dto

import "time"

type VoucherItemResponse struct {
	ID         int64         `json:"id"`
	VoucherID  int64         `json:"voucher_id"`
	ItemID     int64         `json:"item_id"`
	MinQtyDisc float64       `json:"min_qty_disc"`
	CreatedAt  time.Time     `json:"created_at"`
	Item       *ItemResponse `json:"item"`
}

type VoucherItemCreateRequest struct {
	ItemID     int64   `json:"item_id"`
	MinQtyDisc float64 `json:"min_qty_disc"`
}

type VoucherItemRequestGet struct {
	VoucherID int64  `json:"voucher_id"`
	Offset    int64  `json:"offset"`
	Limit     int64  `json:"limit"`
	OrderBy   string `json:"order_by"`
}
