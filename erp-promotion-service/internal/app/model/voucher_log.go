package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type VoucherLog struct {
	ID                    int64     `orm:"column(id);auto" json:"id"`
	VoucherID             int64     `orm:"column(voucher_id)" json:"voucher_id"`
	CustomerID            int64     `orm:"column(customer_id)" json:"customer_id"`
	AddressIDGP           string    `orm:"column(address_id_gp)" json:"address_id_gp"`
	SalesOrderIDGP        string    `orm:"column(sales_order_id_gp)" json:"sales_order_id_gp"`
	VoucherDiscountAmount float64   `orm:"column(vou_disc_amount)" json:"voucher_discount_amount"`
	Status                int8      `orm:"column(status)" json:"status"`
	CreatedAt             time.Time `orm:"column(created_at)" json:"created_at"`
}

func init() {
	orm.RegisterModel(new(VoucherLog))
}

func (m *VoucherLog) TableName() string {
	return "voucher_log"
}

func (m *VoucherLog) MarshalJSON() ([]byte, error) {
	type Alias VoucherLog

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
