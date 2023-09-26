// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID                      int64    `json:"-" valid:"required"`
	PicName                 string   `json:"pic_name" valid:"required"`
	AltPhoneNumber          string   `json:"alt_phone_number"`
	Email                   string   `json:"email"`
	Note                    string   `json:"note"`
	InvoiceTermID           string   `json:"term_invoice_sls_id" valid:"required"`
	PaymentTermID           string   `json:"term_payment_sls_id" valid:"required"`
	PaymentGroupId          string   `json:"payment_group_id" valid:"required"`
	BillingAddress          string   `json:"billing_address" valid:"required"`
	BusinessTypeCreditLimit int8     `json:"business_type_credit_limit"  valid:"required|gt:0"`
	CustomCreditLimit       int8     `json:"custom_credit_limit"`
	CreditLimitAmount       float64  `json:"credit_limit_amount"`
	CreditLimitBefore       float64  `json:"-"`
	CreditLimitAfter        float64  `json:"-"`
	IsCreateCreditLimitLog  int64    `json:"-"`
	KTPPhotos               []string `json:"ktp_photos"`
	MerchantPhotos          []string `json:"merchant_photos"`
	KTPPhotosStr            string   `json:"-"`
	UpdateKTP               bool     `json:"update_ktp"`
	MerchantPhotosStr       string   `json:"-"`

	Merchant         *model.Merchant         `json:"-"`
	PaymentGroup     *model.PaymentGroup     `json:"-"`
	InvoiceTerm      *model.InvoiceTerm      `json:"-"`
	PaymentTerm      *model.SalesTerm        `json:"-"`
	PaymentGroupComb *model.PaymentGroupComb `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	var e error
	var remainingTotalCharge float64

	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Merchant, e = repository.GetMerchant("id", c.ID); e != nil {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	if c.InvoiceTermID != "" {
		InvoiceTermID, _ := common.Decrypt(c.InvoiceTermID)
		c.InvoiceTerm = &model.InvoiceTerm{ID: InvoiceTermID}
		c.InvoiceTerm.Read()
	}

	if c.PaymentTermID != "" {
		PaymentTermID, _ := common.Decrypt(c.PaymentTermID)
		c.PaymentTerm = &model.SalesTerm{ID: PaymentTermID}
		c.PaymentTerm.Read()
	}

	if paymentGroupId, e := common.Decrypt(c.PaymentGroupId); e != nil {
		o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
	} else {
		if c.PaymentGroup, e = repository.ValidPaymentGroup(paymentGroupId); e != nil {
			o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
		} else {
			if c.PaymentGroup.Status != int8(1) {
				o.Failure("payment_group_id.active", util.ErrorActive("payment group"))
			}
		}
	}

	if c.PaymentTerm != nil && c.InvoiceTerm != nil {
		orSelect.Raw("select * from payment_group_comb pgc where pgc.term_payment_sls_id = ? and pgc.term_invoice_sls_id = ?", c.PaymentTerm.ID, c.InvoiceTerm.ID).QueryRow(&c.PaymentGroupComb)
	}

	if c.PaymentGroupComb != nil {
		if c.PaymentGroupComb.PaymentGroup.ID != c.PaymentGroup.ID {
			o.Failure("payment_group_comb.invalid", "Combination Payment Term & Invoice Term is not true")
		}
	} else {
		o.Failure("payment_group_comb.invalid", "Combination Payment Term & Invoice Term is not true")
	}

	c.CreditLimitBefore = c.Merchant.RemainingCreditLimitAmount
	c.CreditLimitAfter = c.Merchant.RemainingCreditLimitAmount

	if c.Merchant.CreditLimit, e = repository.CheckSingleCreditLimitData(c.Merchant.BusinessType.ID, c.PaymentTerm.ID, c.BusinessTypeCreditLimit); e != nil {
		o.Failure("credit_limit.invalid", util.ErrorInvalidData("credit limit"))
	}

	// Add Validation Credit Limit Custom Minus
	if c.CustomCreditLimit == 1 && c.CreditLimitAmount < 0 {
		o.Failure("credit_limit.invalid", util.ErrorCreditCustomMinus("credit limit"))
	}

	if c.CustomCreditLimit != 1 && c.Merchant.CreditLimit != nil {
		c.CreditLimitAmount = c.Merchant.CreditLimit.AmountCreditLimit
	}

	if c.CreditLimitAfter != 0 || c.Merchant.CreditLimitAmount > 0 || c.CreditLimitAmount > 0 {
		if remainingTotalCharge, e = repository.GetTotalRemainingBySOAndSI(c.Merchant.ID); e != nil {
			o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
		}

		c.CreditLimitAfter = c.CreditLimitAmount - remainingTotalCharge
		c.Merchant.RemainingCreditLimitAmount = c.CreditLimitAfter
		c.IsCreateCreditLimitLog = 1
	}

	if len(c.KTPPhotos) > 1 {
		o.Failure("ktp_photos.invalid", util.ErrorEqualLess("photo", "1 KTP Photo"))
	}

	if len(c.MerchantPhotos) > 3 {
		o.Failure("merchant_photos.invalid", util.ErrorEqualLess("photo", "3 Merchant Photo"))
	}

	if len(c.KTPPhotos) > 0 {
		if c.UpdateKTP {
			c.KTPPhotosStr = strings.Join(c.KTPPhotos, ",")
		} else {
			c.KTPPhotosStr = c.Merchant.KTPPhotosUrl
		}
	}

	if len(c.MerchantPhotos) > 0 {
		c.MerchantPhotosStr = strings.Join(c.MerchantPhotos, ",")
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"pic_name.required":                   util.ErrorInputRequired("PIC Name"),
		"term_invoice_sls_id.required":        util.ErrorInputRequired("Term Invoice"),
		"term_payment_sls_id.required":        util.ErrorInputRequired("Term Payment"),
		"billing_address.required":            util.ErrorInputRequired("Billing Address"),
		"payment_group_id.required":           util.ErrorInputRequired("Payment Group"),
		"business_type_credit_limit.gt":       util.ErrorInvalidData("Business Type Credit Limit"),
		"business_type_credit_limit.required": util.ErrorInputRequired("Business Type Credit Limit"),
	}
}
