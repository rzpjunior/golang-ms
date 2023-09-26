// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetSkuDiscount find a single data using field and value condition.
func GetSkuDiscount(field string, values ...interface{}) (*model.SkuDiscount, error) {
	m := new(model.SkuDiscount)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "SkuDiscountItems", 2)
	for i, _ := range m.SkuDiscountItems {
		o.LoadRelated(m.SkuDiscountItems[i], "SkuDiscountItemTiers", 2)
	}

	var quotMark string

	priceSetArr := strings.Split(m.PriceSets, ",")
	for range priceSetArr {
		quotMark += "?,"
	}
	quotMark = strings.TrimSuffix(quotMark, ",")

	o.Raw("select group_concat(name) from price_set where id in ("+quotMark+")", priceSetArr).QueryRow(&m.PriceSetsName)

	quotMark = ""
	orderChannelArr := strings.Split(m.OrderChannels, ",")
	for range orderChannelArr {
		quotMark += "?,"
	}
	quotMark = strings.TrimSuffix(quotMark, ",")

	o.Raw("select group_concat(note) from glossary where attribute = 'order_channel' and value_int in ("+quotMark+")", orderChannelArr).QueryRow(&m.OrderChannelsName)

	return m, nil
}

// GetSkuDiscounts : function to get data from database based on parameters
func GetSkuDiscounts(rq *orm.RequestQuery, priceSetsID string) (sd []*model.SkuDiscount, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SkuDiscount))

	// setting condition to filter by price_set
	cond := q.GetCond()
	if priceSetsID != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("price_set__icontains", ","+priceSetsID+",").Or("price_set__istartswith", priceSetsID+",").Or("price_set__iendswith", ","+priceSetsID).Or("price_set", priceSetsID)

		cond = cond.AndCond(cond1)
	}
	q = q.SetCond(cond)

	if total, err = q.Exclude("status", 3).RelatedSel().All(&sd, rq.Fields...); err != nil {
		return nil, 0, err
	}

	o := orm.NewOrm()
	o.Using("read_only")

	for _, v := range sd {
		var quotMark string
		if v.PriceSets != "" {
			priceSetArr := strings.Split(v.PriceSets, ",")
			for _, _ = range priceSetArr {
				quotMark += "?,"
			}
			quotMark = strings.TrimSuffix(quotMark, ",")

			o.Raw("select group_concat(name) from price_set where id in ("+quotMark+")", priceSetArr).QueryRow(&v.PriceSetsName)

			v.PriceSets = util.EncIdInStr(v.PriceSets)
		}
	}

	return sd, total, nil
}

// ValidSkuDiscount : function to check if data is valid
func ValidSkuDiscount(id int64) (sd *model.SkuDiscount, err error) {
	sd = &model.SkuDiscount{ID: id}
	err = sd.Read("ID")

	return
}

// GetSkuDiscountData : function to get sku discount data based on several parameters (return sku_discount_item data)
func GetSkuDiscountData(merchantID, priceSetID, productID, salesOrderItemID int64, orderChannel int8, date time.Time) (skuDiscountItem *model.SkuDiscountItem, err error) {
	var (
		totalQuotaPerUser, totalDailyQuotaPerUser float64
	)

	o := orm.NewOrm()
	o.Using("read_only")

	q := "select sdi.* " +
		"from sku_discount sd " +
		"join sku_discount_item sdi on sd.id = sdi.sku_discount_id " +
		"where sd.status = 1 " +
		"and find_in_set(?, sd.price_set) " +
		"and ? between sd.start_timestamp and sd.end_timestamp " +
		"and find_in_set(?, sd.order_channel) " +
		"and ((sdi.use_budget = 1 and sdi.rem_budget > 0) or (sdi.use_budget = 2)) " +
		"and sdi.product_id = ? " +
		"order by sd.id asc " +
		"limit 1"
	if err = o.Raw(q, priceSetID, date.Format("2006-01-02 15:04:05"), orderChannel, productID).QueryRow(&skuDiscountItem); err != nil {
		return nil, nil
	}

	totalQuotaPerUser, totalDailyQuotaPerUser, err = GetUsedSkuDiscountData(skuDiscountItem.ID, merchantID, salesOrderItemID)

	salesOrderItem := &model.SalesOrderItem{ID: salesOrderItemID}
	if err = salesOrderItem.Read("ID"); err != nil {
		salesOrderItem = nil
	}
	if salesOrderItem != nil {
		skuDiscountItem.RemOverallQuota = skuDiscountItem.RemOverallQuota + salesOrderItem.DiscountQty
		skuDiscountItem.RemBudget = skuDiscountItem.RemOverallQuota + (salesOrderItem.UnitPriceDiscount * salesOrderItem.DiscountQty)
	}
	skuDiscountItem.RemQuotaPerUser = skuDiscountItem.OverallQuotaPerUser - int64(totalQuotaPerUser)
	skuDiscountItem.RemDailyQuotaPerUser = skuDiscountItem.DailyQuotaPerUser - int64(totalDailyQuotaPerUser)

	skuDiscountItem.SkuDiscount.Read("ID")
	o.LoadRelated(skuDiscountItem, "SkuDiscountItemTiers", 0)

	return skuDiscountItem, nil
}

// GetUsedSkuDiscountData : func to get sku discount data that has been used
func GetUsedSkuDiscountData(skuDiscountItemID, merchantID, salesOrderItemID int64) (totalQuotaPerUser, totalDailyQuotaPerUser float64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var (
		whereSOI       string
		skuDiscountLog []*model.SkuDiscountLog
	)
	if salesOrderItemID != 0 {
		whereSOI = " and sdl.sales_order_item_id != ?"
	} else {
		whereSOI = " and ?=0"
	}

	if _, err = o.Raw("select sdl.* from sku_discount_log sdl where sdl.status = 1 and sdl.sku_discount_item_id = ? and sdl.merchant_id = ?"+whereSOI, skuDiscountItemID, merchantID, salesOrderItemID).QueryRows(&skuDiscountLog); err == nil {
		currentTime := time.Now()
		for _, v := range skuDiscountLog {
			totalQuotaPerUser += v.DiscountQty
			if v.CreatedAt.Format("2006-01-02") == currentTime.Format("2006-01-02") {
				totalDailyQuotaPerUser += v.DiscountQty
			}
		}
	}

	return
}
