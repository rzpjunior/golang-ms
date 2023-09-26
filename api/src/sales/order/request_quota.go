// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// quotaRequest : struct to hold remaining quota data
type quotaRequest struct {
	CurrentTime  time.Time
	Branch       *model.Branch
	Products     []*salesOrderItem
	OrderChannel int8
}

// Validate : function to validate data
func (r *quotaRequest) Validate() *validation.Output {
	var err error
	o := &validation.Output{Valid: true}

	for i, v := range r.Products {
		var maxRemQuota float64

		if v.SkuDiscountItem, err = repository.GetSkuDiscountData(r.Branch.Merchant.ID, r.Branch.PriceSet.ID, v.Product.ID, 0, r.OrderChannel, r.CurrentTime); err == nil && v.SkuDiscountItem != nil {
			// set maximum available qty for discount
			maxRemQuota = float64(v.SkuDiscountItem.RemOverallQuota)
			if float64(v.SkuDiscountItem.RemQuotaPerUser) < maxRemQuota {
				maxRemQuota = float64(v.SkuDiscountItem.RemQuotaPerUser)
			}

			if float64(v.SkuDiscountItem.RemDailyQuotaPerUser) < maxRemQuota {
				maxRemQuota = float64(v.SkuDiscountItem.RemDailyQuotaPerUser)
			}

			if v.SkuDiscountItem.RemOverallQuota <= 0 || (v.SkuDiscountItem.IsUseBudget == 1 && v.SkuDiscountItem.RemBudget <= 0) {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorRunOut("Discount quota for this product"))
			}

			if maxRemQuota > 0 && maxRemQuota != v.MaxDiscQty && maxRemQuota < v.DiscQty {
				o.Failure("rem_qty"+strconv.Itoa(i)+".invalid", util.ErrorHasChange("Max discount quota"))
			}
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *quotaRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
