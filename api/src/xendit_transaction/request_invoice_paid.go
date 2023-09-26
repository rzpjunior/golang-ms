// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package xendit_transaction

import (
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type invoicePaidRequest struct {
	VaNumber               string                `json:"account_number" valid:"required"`
	TransactionDate        string                `json:"transaction_date"`
	PaymentChannelValue    string                `json:"payment_channel"`
	TransactionTime        string                `json:"transaction_time"`
	Amount                 float64               `json:"paid_amount"`
	Token                  string                `json:"token"`
	StatusInvoice          string                `json:"status_invoice"`
	FromCronJob            int                   `json:"from_cronjob"`
	TransactionDateAt      time.Time             `json:"-"`
	TransactionTimeAt      time.Time             `json:"-"`
	ExternalID             string                `json:"external_id" valid:"required"`
	SalesOrderID           string                `json:"-"`
	CreditLimitBefore      float64               `json:"-"`
	CreditLimitAfter       float64               `json:"-"`
	IsCreateCreditLimitLog bool                  `json:"-"`
	SalesOrder             *model.SalesOrder     `json:"-"`
	PaymentChannel         *model.PaymentChannel `json:"-"`
}

func (c *invoicePaidRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	//var code string
	key := []byte("joyfuls joy j0y5")
	dec := decrypt(key, c.Token)
	token := fmt.Sprintf("%s", dec)
	if token != "hey please push on thursday, it will make me happy" {
		o.Failure("token", "invalid token")
	}

	if c.TransactionDate != "" {
		if c.TransactionDateAt, e = time.Parse("2006-01-02", c.TransactionDate); e != nil {
			o.Failure("transaction_date.invalid", "invalid date")
		}
	}
	if e = orSelect.Raw("SELECT id FROM sales_order WHERE code = ? ", c.ExternalID).QueryRow(&c.SalesOrderID); e == nil {
		c.SalesOrder, e = repository.GetSalesOrder("id", c.SalesOrderID)
		if e != nil {
			o.Failure("external_id.invalid", e.Error())
		}

		if e = c.SalesOrder.Branch.Read("ID"); e != nil {
			o.Failure("external_id.invalid", e.Error())
		}

		if e = c.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
			o.Failure("external_id.invalid", e.Error())
		}

		if e = c.SalesOrder.Branch.Merchant.UserMerchant.Read("ID"); e != nil {
			o.Failure("external_id.invalid", e.Error())
		}
		// ======================================================Create Invoice ========================================================
		orSelect.Raw("SELECT * FROM payment_channel WHERE value = ?", c.PaymentChannelValue).QueryRow(&c.PaymentChannel)
		// ==============================================================================================================

	} else {
		o.Failure("external_id.invalid", "external_id invalid or no row found")
	}

	if c.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 {
		c.CreditLimitBefore = c.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
		c.CreditLimitAfter = c.CreditLimitBefore + c.SalesOrder.TotalCharge
		c.IsCreateCreditLimitLog = true
	}

	return o
}

func (c *invoicePaidRequest) Messages() map[string]string {
	return map[string]string{}
}
