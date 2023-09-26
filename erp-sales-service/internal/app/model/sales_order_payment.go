// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(SalesOrderPayment))
}

func (m *SalesOrderPayment) TableName() string {
	return "sales_order_payment"
}

// Sales Order Payment: struct to hold model data for database
type SalesOrderPayment struct {
	ID              int64  `orm:"column(id)" json:"id"`
	SalesOrderID    int64  `orm:"column(sales_order_id)" json:"sales_order_id"`
	CashReceiptIdGP string `orm:"column(cash_receipt_id_gp)" json:"cash_receipt_id_gp"`
}
