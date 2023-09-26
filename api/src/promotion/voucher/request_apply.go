// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package voucher

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type applyRequest struct {
	RedeemCode string `json:"redeem_code"`
	BranchID   string `json:"branch_id" valid:"required"`
	AreaID     string `json:"area_id" valid:"required"`

	Items []*items `json:"items" valid:"required"` // Product list

	Voucher            *model.Voucher            `json:"-"`
	Branch             *model.Branch             `json:"-"`
	Area               *model.Area               `json:"-"`
	AreaBusinessPolicy *model.AreaBusinessPolicy `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type items struct {
	ProductID        string  `json:"product_id"`
	Quantity         float64 `json:"qty"`
	UnitPrice        float64 `json:"unit_price"`
	Note             string  `json:"item_note"`
	TotalSKUDiscount float64 `json:"unit_price_discount"`

	Product *model.Product `json:"-"`
}

// Validate : function to validate request data
func (c *applyRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")
	var err error
	var totalMatch, countMatch int
	var voucher *model.Voucher
	currentTime := time.Now()
	totalOrder := float64(0)
	deliveryFee := float64(0)

	c.Voucher = &model.Voucher{RedeemCode: c.RedeemCode, Status: 1}
	if err = c.Voucher.Read("RedeemCode", "Status"); err != nil {

		// check if all vouchers is archived
		q.Raw("SELECT * FROM voucher where redeem_code = ? ORDER BY id DESC LIMIT 1", c.RedeemCode).QueryRow(&voucher)
		if voucher != nil {
			if voucher.Status != 1 {
				o.Failure("redeem_code.inactive", util.ErrorActive("voucher"))
				return o
			}
		} else {
			o.Failure("redeem_code.invalid", util.ErrorNotFound("voucher"))
			return o
		}

	}

	if branchID, err := common.Decrypt(c.BranchID); err == nil {
		if c.Branch, err = repository.ValidBranch(branchID); err == nil {
			if err = c.Branch.Read("ID"); err == nil {
				c.Branch.Merchant.Read("ID")

				if c.Branch.Area, err = repository.ValidArea(c.Branch.Area.ID); err == nil {
					if c.AreaBusinessPolicy, err = repository.GetAreaBusinessPolicyDelivery(c.Branch.Area.ID, c.Branch.Merchant.BusinessType.ID); err != nil {
						o.Failure("area_business_config_id.invalid", util.ErrorInvalidData("area business config"))
						return o
					}

					deliveryFee = c.AreaBusinessPolicy.DeliveryFee
				} else {
					o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
				}
			} else {
				o.Failure("redeem_code.invalid", util.ErrorInvalidData("branch"))
				return o
			}
		}
	} else {
		o.Failure("redeem_code.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	// Validation: voucher suits with the following merchant
	if c.Voucher.MerchantID != 0 {
		if c.Voucher.MerchantID != c.Branch.Merchant.ID {
			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "merchant"))
		}
	}

	for _, v := range c.Items {

		// Validation: voucher items
		if c.Voucher.VoucherItem == 1 {
			o := orm.NewOrm()
			o.Using("read_only")

			productID, _ := common.Decrypt(v.ProductID)

			o.Raw("SELECT COUNT(*) FROM voucher_item where voucher_id = ? AND product_id = ? AND min_qty_disc <= ? ", c.Voucher.ID, productID, v.Quantity).QueryRow(&countMatch)
			if countMatch > 0 {
				totalMatch++
			}
		}
		totalOrder = totalOrder + (v.Quantity * v.UnitPrice) - v.TotalSKUDiscount

	}

	if totalOrder >= c.AreaBusinessPolicy.MinOrder {
		deliveryFee = 0
	}

	grandTotal := totalOrder + deliveryFee

	if c.RedeemCode != "" {
		if c.Voucher.Status != 1 {
			o.Failure("redeem_code.inactive", util.ErrorActive("voucher"))
			return o
		}

		if currentTime.Before(c.Voucher.StartTimestamp) {
			o.Failure("redeem_code.invalid", util.ErrorNotInPeriod("voucher"))
			return o
		}

		if currentTime.After(c.Voucher.EndTimestamp) {
			o.Failure("redeem_code.invalid", util.ErrorOutOfPeriod("voucher"))
			return o
		}

		if c.Voucher.Type == 1 { //type total discount
			if c.Voucher.MinOrder > totalOrder {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "minimum order"))
				return o
			}

			if c.Voucher.DiscAmount > totalOrder {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "discount amount"))
				return o
			}
		} else if c.Voucher.Type == 2 { // type grand total discount
			if c.Voucher.MinOrder > grandTotal {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("grand total order", "minimum order"))
				return o
			}

			if c.Voucher.DiscAmount > grandTotal {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("grand total order", "discount amount"))
				return o
			}
		} else if c.Voucher.Type == 3 { // type delivery discount
			if c.Voucher.MinOrder > totalOrder {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "minimum order"))
				return o
			}

			if c.Voucher.DiscAmount > deliveryFee {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("delivery fee", "discount amount"))
				return o
			}
		}

		if c.Voucher.RemOverallQuota < 1 {
			o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
			return o
		}

		filter := map[string]interface{}{
			"merchant_id": c.Branch.Merchant.ID,
			"voucher_id":  c.Voucher.ID,
			"status":      int8(1),
		}
		exclude := map[string]interface{}{}
		if _, countVoucherLog, err := repository.CheckVoucherLogData(filter, exclude); err == nil && countVoucherLog >= c.Voucher.UserQuota {
			o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
			return o
		}

		if c.Voucher.TagCustomer != "" {
			sameTagCustomer := ""
			for _, v := range strings.Split(c.Branch.Merchant.TagCustomer, ",") {
				if strings.Contains(c.Voucher.TagCustomer, v) {
					sameTagCustomer = sameTagCustomer + "," + v
				}
			}

			sameTagCustomer = strings.Trim(sameTagCustomer, ",")
			if sameTagCustomer == "" {
				o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "customer tag"))
				return o
			}
		}

		if c.Voucher.Area.ID != 1 && c.Branch.Area.ID != c.Voucher.Area.ID {
			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "area"))
			return o
		}

		c.Voucher.Archetype.Read("ID")
		c.Voucher.Archetype.BusinessType.Read("ID")
		if c.Voucher.Archetype.BusinessType.AuxData != 1 {
			if c.Voucher.Archetype.AuxData != 1 {
				if c.Branch.Archetype.ID != c.Voucher.Archetype.ID {
					o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "archetype"))
					return o
				}
			} else {
				c.Branch.Archetype.Read("ID")
				if c.Branch.Archetype.BusinessType.ID != c.Voucher.Archetype.BusinessType.ID {
					o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "business type"))
					return o
				}
			}
		}

		if c.Voucher.VoucherItem == 1 {
			q.LoadRelated(c.Voucher, "VoucherItems", 0)

			if totalMatch != len(c.Voucher.VoucherItems) {
				o.Failure("redeem_code.invalid", util.ErrorNotValidTermConditions())
			}
		}
	}

	return o
}

// Messages : function to return error messages after validation
func (c *applyRequest) Messages() map[string]string {
	return map[string]string{
		"branch_id.required": util.ErrorInputRequired("branch"),
		"area_id.required":   util.ErrorInputRequired("area"),
	}
}
