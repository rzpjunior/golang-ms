// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

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

// createRequest : struct to hold outlet request data
type createRequest struct {
	CodeUserMerchant string
	CodeMerchant     string
	CodeBranch       string

	MerchantID              string   `json:"merchant_id"`
	MerchantName            string   `json:"merchant_name"`
	MerchantPicName         string   `json:"merchant_pic_name"`
	MerchantPhoneNumber     string   `json:"merchant_phone_number"`
	MerchantAltPhoneNumber  string   `json:"merchant_alt_phone_number"`
	MerchantEmail           string   `json:"merchant_email"`
	CustomerTag             []string `json:"customer_tag"`
	CustomerTagStr          string   `json:"-"`
	MerchantNote            string   `json:"merchant_note"`
	ReferrerCode            string   `json:"referrer_code"`
	ReferralCode            string   `json:"-"`
	ReferenceInfo           string   `json:"reference_info"`
	BusinessTypeCreditLimit int8     `json:"business_type_credit_limit"`

	FinanceAreaID  string `json:"finance_area_id"`
	InvoiceTermID  string `json:"invoice_term_id"`
	PaymentTermID  string `json:"sales_term_id"`
	PaymentGroupID string `json:"payment_group_id"`
	BillingAddress string `json:"billing_address"`

	BranchName            string `json:"branch_name" valid:"required"`
	BranchArchetypeID     string `json:"archetype_id" valid:"required"`
	BranchSalesPersonID   string `json:"salesperson_id" valid:"required"`
	BranchPriceSetID      string `json:"price_set_id" valid:"required"`
	BranchPicName         string `json:"branch_pic_name" valid:"required"`
	BranchPhoneNumber     string `json:"branch_phone_number" valid:"required"`
	BranchAltPhoneNumber  string `json:"branch_alt_phone_number"`
	BranchShippingAddress string `json:"shipping_address" valid:"required"`
	BranchNote            string `json:"branch_note"`
	BranchAreaID          string `json:"branch_area_id" valid:"required"`
	SubDistrictID         string `json:"sub_district_id" valid:"required"`
	WarehouseID           string `json:"warehouse_id" valid:"required"`
	ProspectCustomerID    string `json:"prospect_customer_id"`
	NewMerchantCheck      string `json:"new_merchant_check"`

	MerchantCreditLimitAmount float64 `json:"-"`
	MerchantCreditLimitBefore float64 `json:"-"`
	MerchantCreditLimitAfter  float64 `json:"-"`
	IsCreateCreditLimitLog    bool    `json:"-"`

	KTPPhotos         []string `json:"ktp_photos"`
	MerchantPhotos    []string `json:"merchant_photos"`
	KTPPhotosStr      string   `json:"-"`
	MerchantPhotosStr string   `json:"-"`

	MerchantCreditLimit  *model.CreditLimit
	Merchant             *model.Merchant
	MerchantBusinessType *model.BusinessType
	FinanceArea          *model.Area
	InvoiceTerm          *model.InvoiceTerm
	PaymentTerm          *model.SalesTerm
	PaymentGroup         *model.PaymentGroup
	Referrer             *model.Merchant

	BranchArchetype   *model.Archetype
	BranchSalesPerson *model.Staff
	BranchPriceSet    *model.PriceSet
	BranchArea        *model.Area
	SubDistrict       *model.SubDistrict
	Warehouse         *model.Warehouse
	ProspectCustomer  *model.ProspectCustomer
	Branch            *model.Branch

	Session *auth.SessionData
	Staff   *model.Staff
}

// Validate : function to validate customer tag request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	r := orm.NewOrm()
	r.Using("read_only")
	var (
		err                                                                                                error
		arrCustomerTagInt                                                                                  []int
		paymentTermID, invoiceTermID, paymentGroupID, subDistrictID, branchArchetypeID, prospectCustomerID int64
		existPhoneNum                                                                                      bool
		remainingTotalCharge                                                                               float64
	)

	if salespersonId, e := common.Decrypt(c.BranchSalesPersonID); e != nil {
		o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
	} else {
		if c.BranchSalesPerson, e = repository.ValidStaff(salespersonId); e != nil {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		} else {
			if c.BranchSalesPerson.Status != int8(1) {
				o.Failure("salesperson_id.active", util.ErrorActive("salesperson"))
			}
		}
	}

	if branchArchetypeID, err = common.Decrypt(c.BranchArchetypeID); err != nil {
		o.Failure("branch_archetype.invalid", util.ErrorInvalidData("branch archetype"))
	}

	c.BranchArchetype = &model.Archetype{ID: branchArchetypeID}

	if err = c.BranchArchetype.Read("ID"); err != nil {
		o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
	}

	c.MerchantBusinessType = &model.BusinessType{ID: c.BranchArchetype.BusinessType.ID}
	if err = c.MerchantBusinessType.Read("ID"); err != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	if c.ProspectCustomerID != "" {
		if prospectCustomerID, err = common.Decrypt(c.ProspectCustomerID); err != nil {
			o.Failure("prospect_customer_id.invalid", util.ErrorInvalidData("prospect customer"))
		}
		c.ProspectCustomer = &model.ProspectCustomer{ID: prospectCustomerID}
		c.ProspectCustomer.Read("ID")
		if c.ProspectCustomer.SalespersonID != 0 {
			r.Raw("SELECT * FROM staff where id = ?", c.BranchSalesPerson.ID).QueryRow(&c.Staff)
			c.Staff.User.Read("ID")
		}
	}

	if len(c.CustomerTag) > 0 {
		for _, v := range c.CustomerTag {
			customerId, _ := common.Decrypt(v)

			if customerTag, err := repository.ValidCustomerTag(customerId); err != nil {
				o.Failure("customer_tag.invalid", util.ErrorInvalidData("customer tag"))
			} else {
				if customerTag.Status != int8(1) {
					o.Failure("customer_tag.active", util.ErrorActive("customer tag"))
				}

				arrCustomerTagInt = append(arrCustomerTagInt, int(customerId))
			}
		}

		// sort integer decrypted customer tag id, then convert it into a string with comma separator
		sort.Ints(arrCustomerTagInt)
		customerTagJson, _ := json.Marshal(arrCustomerTagInt)
		c.CustomerTagStr = strings.Trim(string(customerTagJson), "[]")
	}

	if c.NewMerchantCheck == "true" {
		if c.CodeUserMerchant, err = util.CheckTable("user_merchant"); err != nil {
			o.Failure("code.invalid", util.ErrorInvalidData("code"))
		}

		if c.CodeMerchant, err = util.CheckTable("merchant"); err != nil {
			o.Failure("code.invalid", util.ErrorInvalidData("code"))
		}

		filter := map[string]interface{}{"name": c.MerchantName}
		exclude := map[string]interface{}{"status": int8(3)}
		if _, countMerchant, err := repository.CheckMerchantData(filter, exclude); err != nil {
			o.Failure("merchant_name.invalid", util.ErrorInvalidData("merchant name"))
		} else if countMerchant > 0 {
			o.Failure("merchant_name.unique", util.ErrorUnique("merchant name"))
		}

		if c.MerchantName != "" {
			if len(c.MerchantName) > 50 {
				o.Failure("merchant_name.invalid", util.ErrorCharLength("merchant name", 50))
			}
		} else {
			o.Failure("merchant_name.required", util.ErrorInputRequired("merchant name"))
		}

		if c.MerchantPicName == "" {
			o.Failure("merchant_pic_name.required", util.ErrorInputRequired("pic name"))
		}

		if c.BusinessTypeCreditLimit == 0 {
			o.Failure("merchant_business_type_credit_limit.required", util.ErrorInputRequired("merchant business type credit limit"))
		}

		if c.MerchantPhoneNumber == "" {
			o.Failure("phone_number.required", util.ErrorInputRequired("phone number"))
		} else {
			if len(c.MerchantPhoneNumber) < 8 {
				o.Failure("phone_number.invalid", util.ErrorCharLength("phone number", 8))
			} else {
				existPhoneNum = r.QueryTable("merchant").Filter("phone_number", strings.TrimPrefix(c.MerchantPhoneNumber, "0")).Filter("status__in", 1, 2).Exist()
				if existPhoneNum {
					o.Failure("phone_number.invalid", util.ErrorPhoneNumber())
				}
			}
		}

		if c.FinanceAreaID == "" {
			o.Failure("finance_area_id.required", util.ErrorInputRequired("finance area"))
		} else {
			financeAreaID, _ := common.Decrypt(c.FinanceAreaID)
			c.FinanceArea = &model.Area{ID: financeAreaID}
			if err = c.FinanceArea.Read("ID"); err != nil {
				o.Failure("finance_area_id.invalid", util.ErrorInvalidData("finance area"))
			}
		}

		if c.InvoiceTermID == "" {
			o.Failure("invoice_term_id.required", util.ErrorInputRequired("invoice term"))
		} else {
			invoiceTermID, _ = common.Decrypt(c.InvoiceTermID)
			c.InvoiceTerm = &model.InvoiceTerm{ID: invoiceTermID}
			if err = c.InvoiceTerm.Read("ID"); err != nil {
				o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
			}
		}

		if c.PaymentTermID == "" {
			o.Failure("payment_term_id.required", util.ErrorInputRequired("payment term"))
		} else {
			paymentTermID, _ = common.Decrypt(c.PaymentTermID)
			c.PaymentTerm = &model.SalesTerm{ID: paymentTermID}
			if err = c.PaymentTerm.Read("ID"); err != nil {
				o.Failure("payment_term_id.invalid", util.ErrorInvalidData("payment term"))
			}
		}

		if c.PaymentGroupID == "" {
			o.Failure("payment_group_id.required", util.ErrorInputRequired("payment group"))
		} else {
			paymentGroupID, _ = common.Decrypt(c.PaymentGroupID)
			c.PaymentGroup = &model.PaymentGroup{ID: paymentGroupID}
			if err = c.PaymentGroup.Read("ID"); err != nil {
				o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
			}
		}

		if c.BillingAddress != "" {
			if len(c.BillingAddress) > 250 {
				o.Failure("billing_address.invalid", util.ErrorCharLength("billing address", 250))
			}
		} else {
			o.Failure("billing_address.required", util.ErrorInputRequired("billing address"))
		}

		if len(c.BranchName) > 60 {
			o.Failure("branch_name.invalid", util.ErrorCharLength("branch name", 60))
		}

		filter = map[string]interface{}{"term_payment_sls_id": paymentTermID, "term_invoice_sls_id": invoiceTermID, "payment_group_sls_id": paymentGroupID}
		exclude = map[string]interface{}{}
		if _, countPaymentGroupComb, err := repository.CheckPaymentGroupCombData(filter, exclude); err != nil || countPaymentGroupComb == 0 {
			o.Failure("payment_group_comb.invalid", util.ErrorPaymentCombination())
		}

		if c.ReferrerCode != "" {
			c.Referrer = &model.Merchant{ReferralCode: c.ReferrerCode, Status: 1}
			if err = c.Referrer.Read("ReferralCode", "Status"); err != nil {
				o.Failure("referrer_code", util.ErrorMustExistInActive("referrer code", "merchant"))
			}
		}

		if c.MerchantCreditLimit, err = repository.CheckSingleCreditLimitData(c.MerchantBusinessType.ID, c.PaymentTerm.ID, c.BusinessTypeCreditLimit); err != nil {
			o.Failure("credit_limit.invalid", util.ErrorInvalidData("credit limit"))
		}

		if c.MerchantCreditLimit != nil {
			c.MerchantCreditLimitBefore = c.MerchantCreditLimit.AmountCreditLimit
			c.MerchantCreditLimitAfter = c.MerchantCreditLimit.AmountCreditLimit
			c.MerchantCreditLimitAmount = c.MerchantCreditLimit.AmountCreditLimit
			c.IsCreateCreditLimitLog = true
		}
	} else {
		if c.MerchantID != "" {
			merchantID, _ := common.Decrypt(c.MerchantID)
			c.Merchant = &model.Merchant{ID: merchantID}
			if err = c.Merchant.Read("ID"); err == nil {
				if c.ProspectCustomerID != "" && c.Merchant.UpgradeStatus == 1 {
					if c.FinanceAreaID == "" {
						o.Failure("finance_area_id.required", util.ErrorInputRequired("finance area"))
					} else {
						financeAreaID, _ := common.Decrypt(c.FinanceAreaID)
						c.FinanceArea = &model.Area{ID: financeAreaID}
						if err = c.FinanceArea.Read("ID"); err != nil {
							o.Failure("finance_area_id.invalid", util.ErrorInvalidData("finance area"))
						}
					}

					if c.InvoiceTermID == "" {
						o.Failure("invoice_term_id.required", util.ErrorInputRequired("invoice term"))
					} else {
						invoiceTermID, _ = common.Decrypt(c.InvoiceTermID)
						c.InvoiceTerm = &model.InvoiceTerm{ID: invoiceTermID}
						if err = c.InvoiceTerm.Read("ID"); err != nil {
							o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
						}
					}

					if c.PaymentTermID == "" {
						o.Failure("payment_term_id.required", util.ErrorInputRequired("payment term"))
					} else {
						paymentTermID, _ = common.Decrypt(c.PaymentTermID)
						c.PaymentTerm = &model.SalesTerm{ID: paymentTermID}
						if err = c.PaymentTerm.Read("ID"); err != nil {
							o.Failure("payment_term_id.invalid", util.ErrorInvalidData("payment term"))
						}
					}

					if c.PaymentGroupID == "" {
						o.Failure("payment_group_id.required", util.ErrorInputRequired("payment group"))
					} else {
						paymentGroupID, _ = common.Decrypt(c.PaymentGroupID)
						c.PaymentGroup = &model.PaymentGroup{ID: paymentGroupID}
						if err = c.PaymentGroup.Read("ID"); err != nil {
							o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
						}
					}

					if c.ReferrerCode != "" {
						// ReferrerCode cannot be same with referral code of merchant
						if c.Merchant.ReferralCode == c.ReferrerCode {
							o.Failure("referrer_code", util.ErrorReferrerCode())
						}

						// ReferrerCode only update if merchant doesn't have referrer code before
						if c.Merchant.ReferrerCode == "" {
							c.Referrer = &model.Merchant{ReferralCode: c.ReferrerCode, Status: 1}
							if err = c.Referrer.Read("ReferralCode", "Status"); err != nil {
								o.Failure("referrer_code", util.ErrorMustExistInActive("referrer code", "merchant"))
							}
						} else {
							// This Validation used to validate if there is user change the default value referrer code of merchant
							if c.Merchant.ReferrerCode != c.ReferrerCode {
								o.Failure("referrer_code", util.ErrorMustBeSame("referrer code", "referrer code of merchant"))
							}
							c.ReferrerCode = c.Merchant.ReferrerCode
							c.Referrer = c.Merchant.Referrer
						}
					}

					if c.MerchantCreditLimit, err = repository.CheckSingleCreditLimitData(c.MerchantBusinessType.ID, c.PaymentTerm.ID, c.Merchant.BusinessTypeCreditLimit); err != nil {
						o.Failure("credit_limit.invalid", util.ErrorInvalidData("credit limit"))
					}

					if c.MerchantCreditLimit != nil {
						c.MerchantCreditLimitAmount = c.MerchantCreditLimit.AmountCreditLimit
					} else {
						c.MerchantCreditLimitAmount = 0
						c.MerchantCreditLimitAfter = 0
					}

					if c.MerchantCreditLimitAmount > 0 || c.Merchant.CreditLimitAmount > 0 {
						if remainingTotalCharge, err = repository.GetTotalRemainingBySOAndSI(c.Merchant.ID); err != nil {
							o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
						}

						c.MerchantCreditLimitAfter = c.MerchantCreditLimitAmount - remainingTotalCharge
						c.IsCreateCreditLimitLog = true
					}

					c.Merchant.Name = c.MerchantName
					c.Merchant.PicName = c.MerchantPicName
					c.Merchant.PhoneNumber = strings.TrimPrefix(c.MerchantPhoneNumber, "0")
					c.Merchant.AltPhoneNumber = c.MerchantAltPhoneNumber
					c.Merchant.Email = c.MerchantEmail
					c.Merchant.TagCustomer = c.CustomerTagStr
					c.Merchant.Note = c.MerchantNote
					c.Merchant.FinanceArea = c.FinanceArea
					c.Merchant.InvoiceTerm = c.InvoiceTerm
					c.Merchant.PaymentTerm = c.PaymentTerm
					c.Merchant.PaymentGroup = c.PaymentGroup
					c.Merchant.BillingAddress = c.BillingAddress
					c.Merchant.CustomerGroup = 1
					c.Merchant.ReferrerCode = c.ReferrerCode
					c.Merchant.ReferenceInfo = c.ProspectCustomer.ReferenceInfo
					c.Merchant.Referrer = c.Referrer

				}
			} else {
				o.Failure("merchant_id", util.ErrorInvalidData("merchant"))
			}
		} else {
			o.Failure("merchant_id", util.ErrorInputRequired("merchant"))
		}

	}

	if c.CodeBranch, err = util.CheckTable("branch"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if len(c.BranchName) > 60 {
		o.Failure("branch_name.invalid", util.ErrorCharLength("branch name", 60))
	}

	c.BranchName = c.MerchantName + " " + c.BranchName
	filter := map[string]interface{}{"name": c.BranchName}
	exclude := map[string]interface{}{"status": int8(3)}
	if _, countBranch, err := repository.CheckBranchData(filter, exclude); err != nil {
		o.Failure("branch_name.invalid", util.ErrorInvalidData("branch name"))
	} else if countBranch > 0 {
		o.Failure("branch_name.unique", util.ErrorUnique("branch name"))
	}

	branchAreaID, _ := common.Decrypt(c.BranchAreaID)
	c.BranchArea = &model.Area{ID: branchAreaID}
	if err = c.BranchArea.Read("ID"); err != nil {
		o.Failure("branch_area_id", util.ErrorInvalidData("branch area"))
	}

	confApp := &model.ConfigApp{Attribute: "cust_max_tag"}
	confApp.Read("Attribute")
	configAppValue, _ := strconv.Atoi(confApp.Value)
	if len(c.CustomerTag) > configAppValue {
		o.Failure("customer_tag.invalid", util.ErrorSelectMax(confApp.Value, "tag"))
	}

	if c.BranchSalesPersonID != "" {
		branchSalesPersonID, _ := common.Decrypt(c.BranchSalesPersonID)
		c.BranchSalesPerson = &model.Staff{ID: branchSalesPersonID}
		if err = c.BranchSalesPerson.Read("ID"); err != nil {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		}
	}

	branchPriceSetID, _ := common.Decrypt(c.BranchPriceSetID)
	c.BranchPriceSet = &model.PriceSet{ID: branchPriceSetID}
	if err = c.BranchPriceSet.Read("ID"); err != nil {
		o.Failure("price_set_id.invalid", util.ErrorInvalidData("price set"))
	}

	subDistrictID, _ = common.Decrypt(c.SubDistrictID)
	c.SubDistrict = &model.SubDistrict{ID: subDistrictID}
	if err = c.SubDistrict.Read("ID"); err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	}

	warehouseID, _ := common.Decrypt(c.WarehouseID)
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if err = c.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse.Status != 1 {
			o.Failure("warehouse_id.inactive", util.ErrorActive("default warehouse"))
		}

		filter := map[string]interface{}{"sub_district_id": subDistrictID}
		exclude := map[string]interface{}{}
		if _, countWarehouseCoverage, err := repository.CheckWarehouseCoverageData(filter, exclude); err != nil || countWarehouseCoverage == 0 {
			o.Failure("warehouse_id.invalid", util.ErrorMustBeSame("warehouse sub district", "sub district"))
		}
	}

	if len(c.KTPPhotos) > 1 {
		o.Failure("ktp_photos.invalid", util.ErrorEqualLess("photo", "1 KTP Photo"))
	}

	if len(c.MerchantPhotos) > 3 {
		o.Failure("merchant_photos.invalid", util.ErrorEqualLess("photo", "3 Merchant Photo"))
	}

	if len(c.KTPPhotos) > 0 {
		c.KTPPhotosStr = strings.Join(c.KTPPhotos, ",")
	}

	if len(c.MerchantPhotos) > 0 {
		c.MerchantPhotosStr = strings.Join(c.MerchantPhotos, ",")
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"branch_name.required":         util.ErrorInputRequired("branch name"),
		"archetype_id.required":        util.ErrorInputRequired("archetype"),
		"price_set_id.required":        util.ErrorInputRequired("price set"),
		"branch_pic_name.required":     util.ErrorInputRequired("branch pic name"),
		"branch_phone_number.required": util.ErrorInputRequired("branch phone number"),
		"shipping_address.required":    util.ErrorInputRequired("shipping address"),
		"branch_area_id.required":      util.ErrorInputRequired("branch area"),
		"sub_district_id.required":     util.ErrorInputRequired("sub district"),
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
		"salesperson_id.required":      util.ErrorInputRequired("salesperson"),
	}
}
