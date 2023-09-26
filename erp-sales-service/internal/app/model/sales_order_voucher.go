// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(SalesOrderVoucher))
}

func (m *SalesOrderVoucher) TableName() string {
	return "sales_order_voucher"
}

// Sales Order Voucher : struct to hold model data for database
type SalesOrderVoucher struct {
	ID           int64     `orm:"column(id)" json:"id"`
	SalesOrderID int64     `orm:"column(sales_order_id)" json:"sales_order_id"`
	VoucherIDGP  string    `orm:"column(voucher_id_gp)" json:"voucher_id_gp"`
	DiscAmount   float64   `orm:"column(disc_amount)" json:"disc_amount"`
	CreatedAt    time.Time `orm:"column(created_at)" json:"created_at"`
	VoucherType  int8      `orm:"column(voucher_type)" json:"voucher_type"`
}
