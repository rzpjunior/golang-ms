// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(CustomerAccNum))
}

// CustomerAccNum model for city table.
type CustomerAccNum struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	AccountNumber string `orm:"column(account_number);size(100);null" json:"account_number"`
	AccountName   string `orm:"column(account_name);size(100);null" json:"account_name"`

	Customer       *Customer       `orm:"-"  json:"customer"`
	PaymentChannel *PaymentChannel `orm:"-"  json:"payment_channel"`
}
