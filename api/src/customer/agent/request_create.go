// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type priceSetArea struct {
	PriceSetId string `json:"price_set_id"`
	AreaId     string `json:"area_id"`

	PriceSet *model.PriceSet
	Area     *model.Area
}

// createRequest : struct to hold price set request data
type createRequest struct {
	Code               string   `json:"-"`
	UserMerchantCode   string   `json:"-"`
	BranchCode         string   `json:"-"`
	ProspectCustomerId string   `json:"prospect_customer_id"`
	Name               string   `json:"name" valid:"required"`
	PhoneNumber        string   `json:"phone_number" valid:"required"`
	AltPhoneNumber     string   `json:"alt_phone_number"`
	Email              string   `json:"email"`
	CustomerTag        []string `json:"tag_customer"`
	CustomerTagStr     string   `json:"-"`
	Note               string   `json:"note"`
	CustomerGroup      string   `json:"customer_group" valid:"required"`
	BillingAddress     string   `json:"billing_address" valid:"required"`
	//AddressName             string         `json:"address_name" valid:"required"`
	PicName                 string         `json:"pic_name" valid:"required"`
	RecipientName           string         `json:"recipient_name" valid:"required"`
	RecipientPhoneNumber    string         `json:"recipient_phone_number" valid:"required"`
	RecipientAltPhoneNumber string         `json:"recipient_alt_phone_number"`
	ShippingAddress         string         `json:"shipping_address" valid:"required"`
	ShippingNote            string         `json:"shipping_note"`
	BusinessTypeId          string         `json:"business_type_id" valid:"required"`
	BusinessTypeCreditLimit int8           `json:"business_type_credit_limit" valid:"required|gt:0"`
	FinanceAreaId           string         `json:"finance_area_id" valid:"required"`
	InvoiceTermId           string         `json:"term_invoice_sls_id" valid:"required"`
	PaymentTermId           string         `json:"term_payment_sls_id" valid:"required"`
	ArchetypeId             string         `json:"archetype_id" valid:"required"`
	SalespersonId           string         `json:"salesperson_id" valid:"required"`
	SubDistrictId           string         `json:"sub_district_id" valid:"required"`
	PaymentGroupId          string         `json:"payment_group_id" valid:"required"`
	ReferrerCode            string         `json:"referrer_code"`
	ReferralCode            string         `json:"-"`
	ShippingAreaID          string         `json:"shipping_area_id" valid:"required"`
	PriceSetArea            []priceSetArea `json:"price_set_area" valid:"required"`
	CreditLimitAmount       float64        `json:"-"`
	CreditLimitBefore       float64        `json:"-"`
	CreditLimitAfter        float64        `json:"-"`
	IsCreateCreditLimitLog  int64          `json:"-"`

	CreditLimit       *model.CreditLimit
	ProspectCust      *model.ProspectCustomer
	Merchant          *model.Merchant
	BusinessType      *model.BusinessType
	FinanceArea       *model.Area
	InvoiceTerm       *model.InvoiceTerm
	PaymentTerm       *model.SalesTerm
	Archetype         *model.Archetype
	Salesperson       *model.Staff
	PaymentGroup      *model.PaymentGroup
	SubDistrict       *model.SubDistrict
	WarehouseCoverage *model.WarehouseCoverage
	PaymentGroupComb  *model.PaymentGroupComb
	ShippingArea      *model.Area

	Session *auth.SessionData `json:"-"`
	Staff   *model.Staff      `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	r := orm.NewOrm()
	r.Using("read_only")
	var err error
	var arrCustomerTagInt []int
	var existPhoneNum bool

	if c.Code, err = util.CheckTable("merchant"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if c.UserMerchantCode, err = util.CheckTable("user_merchant"); err != nil {
		o.Failure("user_merchant_code.invalid", util.ErrorInvalidData("user merchant code"))
	}

	if c.BranchCode, err = util.CheckTable("branch"); err != nil {
		o.Failure("branch_code.invalid", util.ErrorInvalidData("branch code"))
	}

	if businessTypeId, e := common.Decrypt(c.BusinessTypeId); e != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	} else {
		if c.BusinessType, e = repository.ValidBusinessType(businessTypeId); e != nil {
			o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
		} else {
			if c.BusinessType.Status != int8(1) {
				o.Failure("business_type_id.active", util.ErrorActive("business type"))
			}
		}
	}

	existPhoneNum = r.QueryTable("merchant").Filter("phone_number", strings.TrimPrefix(c.PhoneNumber, "0")).Filter("status__in", 1, 2).Exist()
	if existPhoneNum {
		o.Failure("phone_number.invalid", util.ErrorPhoneNumber())
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number.invalid", util.ErrorCharLength("phone number", 8))
	}

	if c.ReferrerCode != "" {
		var m *model.Merchant
		r.Raw("select * from merchant m where m.referral_code = ?", c.ReferrerCode).QueryRow(&m)
		if m == nil {
			o.Failure("referrer_code.invalid", util.ErrorMustExistInActive("referrer code", "merchant"))
		} else {
			if m.Status != 1 {
				o.Failure("id.invalid", util.ErrorActive("merchant"))
			}
			if len(c.ReferrerCode) < 0 || len(c.ReferrerCode) > 15 {
				o.Failure("id.invalid", "Length of referral code less than 0 or exceed 15")
			}
			c.Merchant = m
		}
	}

	if salespersonId, e := common.Decrypt(c.SalespersonId); e != nil {
		o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
	} else {
		if c.Salesperson, e = repository.ValidStaff(salespersonId); e != nil {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		} else {
			if c.Salesperson.Status != int8(1) {
				o.Failure("salesperson_id.active", util.ErrorActive("salesperson"))
			}
		}
	}

	if c.ProspectCustomerId != "" {
		if prosCustId, e := common.Decrypt(c.ProspectCustomerId); e != nil {
			o.Failure("prospect_customer_id.invalid", util.ErrorInvalidData("prospect customer"))
		} else {
			if c.ProspectCust, e = repository.ValidProspectiveCustomer(prosCustId); e != nil {
				o.Failure("prospect_customer_id.invalid", util.ErrorInvalidData("prospect customer"))
			}
			if c.ProspectCust.SalespersonID != 0 {
				r.Raw("SELECT * FROM staff where id = ?", c.Salesperson.ID).QueryRow(&c.Staff)
				c.Staff.User.Read("ID")
			}
		}
	}

	if financeAreaId, e := common.Decrypt(c.FinanceAreaId); e != nil {
		o.Failure("finance_area_id.invalid", util.ErrorInvalidData("finance area"))
	} else {
		if c.FinanceArea, e = repository.ValidArea(financeAreaId); e != nil {
			o.Failure("finance_area_id.invalid", util.ErrorInvalidData("finance area"))
		} else {
			if c.FinanceArea.Status != int8(1) {
				o.Failure("finance_area_id.active", util.ErrorActive("finance area"))
			}
		}
	}

	if shippingAreaId, e := common.Decrypt(c.ShippingAreaID); e != nil {
		o.Failure("shipping_area_id.invalid", util.ErrorInvalidData("shipping area"))
	} else {
		if c.ShippingArea, e = repository.ValidArea(shippingAreaId); e != nil {
			o.Failure("shipping_area_id.invalid", util.ErrorInvalidData("shipping area"))
		} else {
			if c.ShippingArea.Status != int8(1) {
				o.Failure("shipping_area_id.active", util.ErrorActive("shipping area"))
			}
		}
	}

	if invoiceTermId, e := common.Decrypt(c.InvoiceTermId); e != nil {
		o.Failure("term_invoice_sls_id.invalid", util.ErrorInvalidData("term invoice sales"))
	} else {
		if c.InvoiceTerm, e = repository.ValidInvoiceTerm(invoiceTermId); e != nil {
			o.Failure("term_invoice_sls_id.invalid", util.ErrorInvalidData("term invoice sales"))
		} else {
			if c.InvoiceTerm.Status != int8(1) {
				o.Failure("term_invoice_sls_id.active", util.ErrorActive("term invoice sales"))
			}
		}
	}

	if paymentTermId, e := common.Decrypt(c.PaymentTermId); e != nil {
		o.Failure("term_payment_sls_id.invalid", util.ErrorInvalidData("payment term"))
	} else {
		if c.PaymentTerm, e = repository.ValidSalesTerm(paymentTermId); e != nil {
			o.Failure("term_payment_sls_id.invalid", util.ErrorInvalidData("payment term"))
		} else {
			if c.PaymentTerm.Status != int8(1) {
				o.Failure("term_payment_sls_id.active", util.ErrorActive("payment term"))
			}
		}
	}

	if archetypeId, e := common.Decrypt(c.ArchetypeId); e != nil {
		o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
	} else {
		if c.Archetype, e = repository.ValidArchetype(archetypeId); e != nil {
			o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
		} else {
			if c.Archetype.Status != int8(1) {
				o.Failure("archetype_id.active", util.ErrorActive("archetype"))
			}
		}
	}

	if salespersonId, e := common.Decrypt(c.SalespersonId); e != nil {
		o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
	} else {
		if c.Salesperson, e = repository.ValidStaff(salespersonId); e != nil {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		} else {
			if c.Salesperson.Status != int8(1) {
				o.Failure("salesperson_id.active", util.ErrorActive("salesperson"))
			}
		}
	}

	if len(c.BillingAddress) > 250 {
		o.Failure("billing_address.invalid", util.ErrorCharLength("billing address", 250))
	}

	if len(c.ShippingAddress) > 350 {
		o.Failure("shipping_address.invalid", util.ErrorCharLength("shipping address", 350))
	}

	if subDistrictId, e := common.Decrypt(c.SubDistrictId); e != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	} else {
		if c.SubDistrict, e = repository.ValidSubDistrict(subDistrictId); e != nil {
			o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
		} else {
			if c.SubDistrict.Status != int8(1) {
				o.Failure("sub_district_id.active", util.ErrorActive("sub district"))
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

	confApp := &model.ConfigApp{
		Attribute: "cust_max_tag",
	}
	confApp.Read("Attribute")
	configAppValue, _ := strconv.Atoi(confApp.Value)
	if len(c.CustomerTag) > configAppValue {
		o.Failure("tag_customer.invalid", util.ErrorSelectMax("3", "customer tag"))
	} else {
		if len(c.CustomerTag) > 0 {
			for _, v := range c.CustomerTag {
				customerId, _ := common.Decrypt(v)

				if customerTag, err := repository.ValidCustomerTag(customerId); err != nil {
					o.Failure("tag_customer.invalid", util.ErrorInvalidData("customer tag"))
				} else {
					if customerTag.Status != int8(1) {
						o.Failure("tag_customer.active", util.ErrorActive("customer tag"))
					}

					arrCustomerTagInt = append(arrCustomerTagInt, int(customerId))
				}
			}

			// sort integer decrypted customer tag id, then convert it into a string with comma separator
			sort.Ints(arrCustomerTagInt)
			customerTagJson, _ := json.Marshal(arrCustomerTagInt)
			c.CustomerTagStr = strings.Trim(string(customerTagJson), "[]")
		}
	}

	for k := range c.PriceSetArea {
		if c.PriceSetArea[k].PriceSetId == "" {
			o.Failure("price_set_area"+strconv.Itoa(k), util.ErrorInputRequired("price set"))
		}
		if areaId, e := common.Decrypt(c.PriceSetArea[k].AreaId); e != nil {
			o.Failure("area_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("area"))
		} else {
			if c.PriceSetArea[k].Area, e = repository.ValidArea(areaId); e != nil {
				o.Failure("area_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("area"))
			} else {
				if c.PriceSetArea[k].Area.Status != int8(1) {
					o.Failure("area_id"+strconv.Itoa(k)+".active", util.ErrorActive("area"))
				}
			}
		}

		if c.PriceSetArea[k].AreaId == "" {
			o.Failure("price_set_area"+strconv.Itoa(k)+".required", util.ErrorInputRequired("area"))
		}

		if c.PriceSetArea[k].PriceSetId != "" {
			if priceSetId, e := common.Decrypt(c.PriceSetArea[k].PriceSetId); e != nil {
				o.Failure("price_set_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("price_set"))
			} else {
				if c.PriceSetArea[k].PriceSet, e = repository.ValidPriceSet(priceSetId); e != nil {
					o.Failure("price_set_id"+strconv.Itoa(k)+".invalid", util.ErrorInvalidData("price_set"))
				} else {
					if c.PriceSetArea[k].PriceSet.Status != int8(1) {
						o.Failure("price_set_id"+strconv.Itoa(k)+".active", util.ErrorActive("price_set"))
					}
				}
			}
		}
	}

	if c.SubDistrict != nil && c.PaymentTerm != nil && c.InvoiceTerm != nil {
		r.Raw("select * from warehouse_coverage wc where sub_district_id =? and main_warehouse = 1", c.SubDistrict.ID).QueryRow(&c.WarehouseCoverage)
		r.Raw("select * from payment_group_comb pgc where pgc.term_payment_sls_id = ? and pgc.term_invoice_sls_id = ?", c.PaymentTerm.ID, c.InvoiceTerm.ID).QueryRow(&c.PaymentGroupComb)
	}

	if c.WarehouseCoverage == nil {
		o.Failure("warehouse_coverage.invalid", util.ErrorMustBeSame("warehouse_coverage", "warehouse"))
	}
	if c.PaymentGroupComb != nil {
		if c.PaymentGroupComb.PaymentGroup.ID != c.PaymentGroup.ID {
			o.Failure("payment_group_comb.invalid", "Combination Payment Term & Invoice Term is not true")
		}
	} else {
		o.Failure("payment_group_comb.invalid", "Combination Payment Term & Invoice Term is not true")
	}

	c.ReferralCode = util.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)

	if c.ProspectCustomerId != "" {
		r.Raw("SELECT * FROM staff where id = ?", c.Salesperson.ID).QueryRow(&c.Staff)
		c.Staff.User.Read("ID")
	}

	if c.CreditLimit, err = repository.CheckSingleCreditLimitData(c.BusinessType.ID, c.PaymentTerm.ID, c.BusinessTypeCreditLimit); err != nil {
		o.Failure("credit_limit.invalid", util.ErrorInvalidData("credit limit"))
	}

	if c.CreditLimit != nil {
		c.CreditLimitBefore = c.CreditLimit.AmountCreditLimit
		c.CreditLimitAfter = c.CreditLimit.AmountCreditLimit
		c.CreditLimitAmount = c.CreditLimit.AmountCreditLimit
		c.IsCreateCreditLimitLog = 1
	}

	return o

}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"business_type_id.required":           util.ErrorInputRequired("business type"),
		"business_type_credit_limit.gt":       util.ErrorInvalidData("business type credit limit"),
		"business_type_credit_limit.required": util.ErrorInputRequired("business type credit limit"),
		"customer_group.required":             util.ErrorInputRequired("customer group"),
		"name.required":                       util.ErrorInputRequired("name"),
		"finance_area_id.required":            util.ErrorInputRequired("finance area"),
		"term_invoice_sls_id.required":        util.ErrorInputRequired("invoice term"),
		"term_payment_sls_id.required":        util.ErrorInputRequired("payment term"),
		"payment_method_id.required":          util.ErrorInputRequired("payment method"),
		"billing_address.required":            util.ErrorInputRequired("billing address"),
		"archetype_id.required":               util.ErrorInputRequired("archetype"),
		"price_set_id.required":               util.ErrorInputRequired("price set"),
		"recipient_name.required":             util.ErrorInputRequired("recipient name"),
		"pic_name.required":                   util.ErrorInputRequired("pic name"),
		"recipient_phone_number.required":     util.ErrorInputRequired("recipient phone number"),
		"shipping_address.required":           util.ErrorInputRequired("shipping address"),
		"sub_district_id.required":            util.ErrorInputRequired("sub district"),
		"warehouse_id.required":               util.ErrorInputRequired("warehouse"),
		"payment_group_id.required":           util.ErrorInputRequired("payment group"),
		"phone_number.required":               util.ErrorInputRequired("phone number"),
		"salesperson_id.required":             util.ErrorInputRequired("salesperson"),
		"shipping_area_id.required":           util.ErrorInputRequired("shipping area"),
	}

	for k, _ := range c.PriceSetArea {
		messages["price_set_area."+strconv.Itoa(k)+".area_id.required"] = util.ErrorInputRequired("area")
		messages["price_set_area."+strconv.Itoa(k)+".price_set_id.required"] = util.ErrorInputRequired("price set")
	}

	return messages
}
