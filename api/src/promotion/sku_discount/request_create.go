// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sku_discount

import (
	"strconv"
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

// createRequest : struct to hold request data
type createRequest struct {
	Code              string
	Name              string   `json:"name" valid:"required"`
	PriceSetArr       []string `json:"price_set" valid:"required"`
	StartTimestampStr string   `json:"start_period" valid:"required"`
	EndTimestampStr   string   `json:"end_period" valid:"required"`
	DivisionID        string   `json:"division_id" valid:"required"`
	OrderChannelArr   []string `json:"order_channel" valid:"required"`
	Note              string   `json:"note"`

	PriceSet       string    `json:"-"`
	OrderChannel   string    `json:"-"`
	StartTimestamp time.Time `json:"-"`
	EndTimestamp   time.Time `json:"-"`

	Division *model.Division `json:"-"`
	Items    []*requestItems `json:"items"`

	Session *auth.SessionData
}

type requestItems struct {
	ProductID           string       `json:"product_id"`
	OverallQuota        int64        `json:"overall_quota"`
	OverallQuotaPerUser int64        `json:"overall_quota_per_user"`
	DailyQuotaPerUser   int64        `json:"daily_quota_per_user"`
	Budget              float64      `json:"budget"`
	IsUseBudget         int8         `json:"-"`
	Tiers               []*itemTiers `json:"item_tier"`

	Product *model.Product `json:"-"`
}

type itemTiers struct {
	TierLevel  int8    `json:"-"`
	MinimumQty float64 `json:"minimum_qty"`
	Amount     float64 `json:"amount"`
}

// Validate : function to validate request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                           error
		divisionID, productID         int64
		valuesPS, valuesPrd           []int64
		valuesOC                      []string
		q, wherePS, whereOC, wherePrd string
		isDataExist                   bool
		currentTime                   time.Time
	)

	orm := orm.NewOrm()
	orm.Using("read_only")
	productList := make(map[int64]string)
	wib, _ := time.LoadLocation("Asia/Jakarta")

	if r.Code, err = util.CheckTable("sku_discount"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	productPromotion := &model.SkuDiscount{Name: r.Name, Status: 1}
	if err = productPromotion.Read("Name", "Status"); err == nil {
		o.Failure("name.invalid", util.ErrorDuplicate("name"))
	}

	wherePS = "("
	for _, v := range r.PriceSetArr {
		var (
			priceSetID int64
			priceSet   *model.PriceSet
		)

		if priceSetID, err = common.Decrypt(v); err != nil {
			o.Failure("price_set.invalid", util.ErrorInvalidData("price set"))
			break
		}

		if priceSet, err = repository.ValidPriceSet(priceSetID); err != nil {
			o.Failure("price_set.invalid", util.ErrorInvalidData("price set"))
			break
		}

		if priceSet.Status != 1 {
			o.Failure("price_set.invalid", util.ErrorActive("price set"))
			break
		}

		r.PriceSet += strconv.Itoa(int(priceSetID)) + ","

		wherePS += "find_in_set(?, sd.price_set) or "
		valuesPS = append(valuesPS, priceSetID)
	}
	wherePS = strings.TrimSuffix(wherePS, " or ") + ")"
	r.PriceSet = strings.TrimSuffix(r.PriceSet, ",")

	whereOC = "("
	for _, v := range r.OrderChannelArr {
		// if order channel is not mobile (2 or 3) then return error
		// commented temporarily
		// if !(v == "2" || v == "3") {
		// 	o.Failure("order_channel.invalid", util.ErrorInvalidData("order channel"))
		// }

		r.OrderChannel += v + ","

		whereOC += "find_in_set(?, sd.order_channel) or "
		valuesOC = append(valuesOC, v)
	}
	whereOC = strings.TrimSuffix(whereOC, " or ") + ")"
	r.OrderChannel = strings.TrimSuffix(r.OrderChannel, ",")

	if r.StartTimestamp, err = time.ParseInLocation("2006-01-02 15:04:05", r.StartTimestampStr, wib); err != nil {
		o.Failure("start_period.invalid", util.ErrorInvalidData("start time"))
	} else {
		currentTime, err = time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), wib)
		if r.StartTimestamp.Before(currentTime) {
			o.Failure("start_period.invalid", util.ErrorLater("Start time", "current time"))
		}
	}

	if r.EndTimestamp, err = time.ParseInLocation("2006-01-02 15:04:05", r.EndTimestampStr, wib); err != nil {
		o.Failure("end_period.invalid", util.ErrorInvalidData("end time"))
	}

	if !r.StartTimestamp.IsZero() && !r.EndTimestamp.IsZero() && (r.EndTimestamp.Before(r.StartTimestamp) || r.EndTimestamp.Equal(r.StartTimestamp)) {
		o.Failure("start_period.invalid", util.ErrorLater("start time", "end time"))
	}

	if divisionID, err = common.Decrypt(r.DivisionID); err != nil {
		o.Failure("division_id.invalid", util.ErrorInvalidData("division"))
	} else {
		if r.Division, err = repository.ValidDivision(divisionID); err != nil {
			o.Failure("division_id.invalid", util.ErrorInvalidData("division"))
		} else {
			if r.Division.Status != 1 {
				o.Failure("division_id.invalid", util.ErrorActive("division"))
			}
		}
	}

	for i, v := range r.Items {
		var (
			minimumQty float64
			tierLevel  int8
		)

		if v.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorSelectRequired("product"))
			continue
		}

		if productID, err = common.Decrypt(v.ProductID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			continue
		}

		if v.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			continue
		}

		if _, exist := productList[productID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
			continue
		}

		productList[productID] = "t"

		wherePrd += "?,"
		valuesPrd = append(valuesPrd, v.Product.ID)

		if v.OverallQuota <= 0 {
			o.Failure("overall_quota"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorGreater("", "0"))))
		}

		if v.OverallQuotaPerUser > v.OverallQuota {
			o.Failure("overall_quota_per_user"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorLess("", "overall quota"))))
		}

		if err = v.Product.Uom.Read(("ID")); err != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product uom"))
		}

		v.Tiers = v.Tiers[:len(v.Tiers)-1]
		if len(v.Tiers) > 0 {
			for idx, val := range v.Tiers {
				tierLevel++
				val.TierLevel = tierLevel

				if val.MinimumQty <= minimumQty {
					minLimitStr := "0"
					if minimumQty > 0 {
						minLimitStr = "previous tier"
					}
					o.Failure("minimum_qty"+strconv.Itoa(i)+"_"+strconv.Itoa(idx)+".invalid", strings.Title(strings.TrimSpace(util.ErrorGreater("", minLimitStr))))
				}
				minimumQty = val.MinimumQty

				if idx == 0 && val.Amount <= 0 {
					o.Failure("amount"+strconv.Itoa(i)+"_0.invalid", strings.Title(strings.TrimSpace(util.ErrorGreater("", "0"))))
				}

				if v.Product.Uom.DecimalEnabled == 2 {
					if val.MinimumQty != float64((int64(val.MinimumQty))) {
						o.Failure("minimum_qty"+strconv.Itoa(i)+"_"+strconv.Itoa(idx)+".invalid", util.ErrorInvalidData("product quantity"))
					}
				}
			}
		} else {
			o.Failure("item_tier"+strconv.Itoa(i)+".required", util.ErrorInputRequired("tier data"))
		}

		if v.OverallQuota < int64(minimumQty) {
			o.Failure("overall_quota"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorEqualGreater("", "minimum qty of last tier"))))
		}

		if v.OverallQuotaPerUser > 0 && v.OverallQuotaPerUser < int64(minimumQty) {
			o.Failure("overall_quota_per_user"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorEqualGreater("", "minimum qty of last tier"))))
		} else if v.OverallQuotaPerUser == 0 {
			v.OverallQuotaPerUser = v.OverallQuota
		}

		if v.DailyQuotaPerUser > 0 {
			if v.DailyQuotaPerUser > v.OverallQuotaPerUser {
				o.Failure("daily_quota_per_user"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorLess("", "quota per user"))))
			}

			if v.DailyQuotaPerUser < int64(minimumQty) {
				o.Failure("daily_quota_per_user"+strconv.Itoa(i)+".invalid", strings.Title(strings.TrimSpace(util.ErrorEqualGreater("", "minimum qty of last tier"))))
			}
		} else {
			v.DailyQuotaPerUser = v.OverallQuotaPerUser
		}

		v.IsUseBudget = 2
		if v.Budget > 0 {
			v.IsUseBudget = 1
		}
	}
	wherePrd = strings.TrimSuffix(wherePrd, ",")

	// validate if there are active sku discount that overlaps with one that is going to be created
	if len(valuesPrd) > 0 && len(valuesPS) > 0 && len(valuesOC) > 0 {
		q = "select exists(select sd.id, sdi.id from sku_discount sd join sku_discount_item sdi on sd.id = sdi.sku_discount_id where sd.status = 1 and " + wherePS +
			" and " + whereOC + " and not (sd.end_timestamp < ? or sd.start_timestamp > ?) and sdi.product_id in (" + wherePrd + "))"
		if err = orm.Raw(q, valuesPS, valuesOC, r.StartTimestamp.Format("2006-01-02 15:04:05"), r.EndTimestamp.Format("2006-01-02 15:04:05"), valuesPrd).QueryRow(&isDataExist); err != nil || (err == nil && isDataExist) {
			o.Failure("id.invalid", "There is same SKU Discount with other promotion")
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":          util.ErrorInputRequired("name"),
		"price_set_id.required":  util.ErrorSelectRequired("price set"),
		"start_period.required":  util.ErrorSelectRequired("start period"),
		"end_period.required":    util.ErrorSelectRequired("end date & end time"),
		"division_id.required":   util.ErrorSelectRequired("division"),
		"order_channel.required": util.ErrorSelectRequired("order channel"),
	}
}
