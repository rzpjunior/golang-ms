package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type VoucherItem struct {
	ID         int64     `orm:"column(id)" json:"id"`
	VoucherID  int64     `orm:"column(voucher_id)" json:"voucher_id"`
	ItemID     int64     `orm:"column(item_id)" json:"item_id"`
	MinQtyDisc float64   `orm:"column(min_qty_disc)" json:"min_qty_disc"`
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(VoucherItem))
}

func (m *VoucherItem) TableName() string {
	return "voucher_item"
}

func (m *VoucherItem) MarshalJSON() ([]byte, error) {
	type Alias VoucherItem

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
