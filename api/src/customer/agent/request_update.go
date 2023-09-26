// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"fmt"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID                      int64          `json:"-" valid:"required"`
	PicName                 string         `json:"pic_name" valid:"required"`
	AltPhoneNumber          string         `json:"alt_phone_number"`
	Email                   string         `json:"email"`
	Note                    string         `json:"note"`
	BillingAddress          string         `json:"billing_address" valid:"required"`
	InvoiceTermId           string         `json:"term_invoice_sls_id" valid:"required"`
	PaymentTermId           string         `json:"term_payment_sls_id" valid:"required"`
	SalespersonId           string         `json:"salesperson_id" valid:"required"`
	PaymentGroupId          string         `json:"payment_group_id" valid:"required"`
	PriceSetArea            []priceSetArea `json:"price_set_area" valid:"required"`
	BusinessTypeCreditLimit int8           `json:"business_type_credit_limit" valid:"required|gt:0"`
	CustomCreditLimit       int8           `json:"custom_credit_limit"`
	CreditLimitAmount       float64        `json:"credit_limit_amount"`
	CreditLimitBefore       float64        `json:"-"`
	CreditLimitAfter        float64        `json:"-"`
	IsCreateCreditLimitLog  int64          `json:"-"`
	NotePriceSetChange      string         `json:"-"`

	Merchant      *model.Merchant
	PaymentGroup  *model.PaymentGroup
	InvoiceTerm   *model.InvoiceTerm
	PaymentTerm   *model.SalesTerm
	PaymentMethod *model.PaymentMethod
	PriceSet      *model.PriceSet
	Salesperson   *model.Staff

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {

	var (
		e                    error
		existingPriceSets    []*model.MerchantPriceSet
		isPriceSetChanged    bool
		priceSetChange       string
		remainingTotalCharge float64 = 0
	)
	oldPriceSet := make(map[int64]*model.PriceSet)

	o := &validation.Output{Valid: true}

	if c.Merchant, e = repository.GetMerchant("id", c.ID); e != nil {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	if invoiceTermId, e := common.Decrypt(c.InvoiceTermId); e != nil {
		o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
	} else {
		if c.InvoiceTerm, e = repository.ValidInvoiceTerm(invoiceTermId); e != nil {
			o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
		} else {
			if c.InvoiceTerm.Status != int8(1) {
				o.Failure("invoice_term_id.active", util.ErrorActive("invoice term"))
			}
		}
	}

	if paymentTermId, e := common.Decrypt(c.PaymentTermId); e != nil {
		o.Failure("payment_term_id.invalid", util.ErrorInvalidData("payment term"))
	} else {
		if c.PaymentTerm, e = repository.ValidSalesTerm(paymentTermId); e != nil {
			o.Failure("payment_term_id.invalid", util.ErrorInvalidData("payment term"))
		} else {
			if c.PaymentTerm.Status != int8(1) {
				o.Failure("payment_term_id.active", util.ErrorActive("payment term"))
			}
		}
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

	// Get Existing Price Sets of Agent
	if existingPriceSets, _, e = repository.GetListPriceSetAgent(c.Merchant.ID); e != nil {
		o.Failure("price_set_area.invalid", util.ErrorInvalidData("price set"))
	}

	for _, v := range existingPriceSets {
		if v.PriceSet, e = repository.ValidPriceSet(v.PriceSet.ID); e != nil {
			o.Failure("price_set_id.invalid", util.ErrorInvalidData("price set"))
		}
		oldPriceSet[v.Area.ID] = v.PriceSet
	}

	for k := range c.PriceSetArea {
		if c.PriceSetArea[k].PriceSetId == "" {
			o.Failure("price_set_area"+strconv.Itoa(k), util.ErrorInputRequired("price set"))
		}
		if areaId, e := common.Decrypt(c.PriceSetArea[k].AreaId); e != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.PriceSetArea[k].Area, e = repository.ValidArea(areaId); e != nil {
				o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
			} else {
				if c.PriceSetArea[k].Area.Status != int8(1) {
					o.Failure("area_id.active", util.ErrorActive("area"))
				}
			}
		}

		if c.PriceSetArea[k].AreaId == "" {
			o.Failure("price_set_area"+strconv.Itoa(k)+"area_id.required", util.ErrorInputRequired("area"))
		}
		if priceSetId, e := common.Decrypt(c.PriceSetArea[k].PriceSetId); e != nil {
			o.Failure("price_set_id.invalid", util.ErrorInvalidData("price set"))
		} else {
			if c.PriceSetArea[k].PriceSet, e = repository.ValidPriceSet(priceSetId); e != nil {
				o.Failure("price_set_id.invalid", util.ErrorInvalidData("price set"))
			} else {
				if c.PriceSetArea[k].PriceSet.Status != int8(1) {
					o.Failure("price_set_id.active", util.ErrorActive("price set"))
				}
			}
		}

		// If there is changes price set area, record price set before and price set after
		if oldPriceSet[c.PriceSetArea[k].Area.ID].ID != c.PriceSetArea[k].PriceSet.ID {
			priceSetChange += fmt.Sprintf("%s: Before: %s - After: %s; ", c.PriceSetArea[k].Area.Name, oldPriceSet[c.PriceSetArea[k].Area.ID].Name, c.PriceSetArea[k].PriceSet.Name)
			isPriceSetChanged = true
		}
	}

	if isPriceSetChanged {
		c.NotePriceSetChange = "Priceset Changed | " + priceSetChange
	}

	c.NotePriceSetChange = strings.TrimSuffix(c.NotePriceSetChange, "; ")

	if salesPersonId, e := common.Decrypt(c.SalespersonId); e != nil {
		o.Failure("sales_person.invalid", util.ErrorInvalidData("salesperson"))
	} else {
		c.Salesperson = &model.Staff{ID: salesPersonId}
		c.Salesperson.Read("ID")
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

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"pic_name.required":                   util.ErrorInputRequired("pic name"),
		"business_type_id.required":           util.ErrorInputRequired("business type"),
		"business_type_credit_limit.gt":       util.ErrorInvalidData("business type credit limit"),
		"business_type_credit_limit.required": util.ErrorInputRequired("business type credit limit"),
		"billing_address.required":            util.ErrorInputRequired("billing address"),
		"term_invoice_sls_id.required":        util.ErrorInputRequired("invoice term"),
		"term_payment_sls_id.required":        util.ErrorInputRequired("payment term"),
		"payment_method_id.required":          util.ErrorInputRequired("payment method"),
		"price_set_id.required":               util.ErrorInputRequired("price set"),
		"salesperson_id.required":             util.ErrorInputRequired("salesperson"),
		"payment_group_id.required":           util.ErrorInputRequired("payment group"),
	}
}
