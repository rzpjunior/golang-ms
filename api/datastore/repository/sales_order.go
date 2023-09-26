// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetSalesOrder find a single data sales term using field and value condition.
func GetSalesOrder(field string, values ...interface{}) (*model.SalesOrder, error) {
	var (
		totalQuotaPerUser, totalDailyQuotaPerUser float64
		err                                       error
	)
	m := new(model.SalesOrder)
	o := orm.NewOrm()
	o.Using("read_only")

	if err = o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "SalesOrderItems", 2)
	o.LoadRelated(m, "DeliveryOrder", 2)
	o.LoadRelated(m, "SalesInvoice", 2)

	o.Raw("select note from audit_log where `type` = 'sales_order' and `function` = 'cancel' and ref_id = ?", m.ID).QueryRow(&m.CancellationNote)
	o.Raw("select name from staff where id = ?", m.LockedBy).QueryRow(&m.LockedByName)
	o.Raw("select poa.status from picking_order_assign poa where poa.sales_order_id = ?", m.ID).QueryRow(&m.StatusPickingOrderAssign)

	for _, v := range m.SalesOrderItems {
		o.LoadRelated(v, "SkuDiscountItem")
		if v.SkuDiscountItem != nil {
			totalQuotaPerUser, totalDailyQuotaPerUser, err = GetUsedSkuDiscountData(v.SkuDiscountItem.ID, m.Branch.Merchant.ID, v.ID)
			v.SkuDiscountItem.RemOverallQuota = v.SkuDiscountItem.RemOverallQuota + v.DiscountQty
			v.SkuDiscountItem.RemQuotaPerUser = v.SkuDiscountItem.OverallQuotaPerUser - int64(totalQuotaPerUser)
			v.SkuDiscountItem.RemDailyQuotaPerUser = v.SkuDiscountItem.DailyQuotaPerUser - int64(totalDailyQuotaPerUser)
			if v.SkuDiscountItem.IsUseBudget == 1 {
				v.SkuDiscountItem.RemBudget = v.SkuDiscountItem.RemBudget + (v.DiscountQty * v.UnitPriceDiscount)
			}

			o.LoadRelated(v.SkuDiscountItem, "SkuDiscountItemTiers")
		}
		o.Raw("select doi.* from delivery_order_item doi join delivery_order do on doi.delivery_order_id = do.id where do.status in (1,2,5,6,7) and doi.sales_order_item_id = ?", v.ID).QueryRow(&v.DeliveryOrderItem)
	}

	return m, nil
}

// GetSalesOrders : function to get data from database based on parameters
func GetSalesOrders(rq *orm.RequestQuery) (m []*model.SalesOrder, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.SalesOrder))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesOrder
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			o.Raw("select poa.status from picking_order_assign poa where poa.sales_order_id = ?", v.ID).QueryRow(&v.StatusPickingOrderAssign)
		}
		return mx, total, nil
	}

	return nil, total, err
}

// CheckSalesOrdersData : function to get all sales order data based on filter and exclude parameters
func CheckSalesOrdersData(filter, exclude map[string]interface{}) (m []*model.SalesOrder, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	qtb := o.QueryTable(new(model.SalesOrder))

	for k, v := range filter {
		qtb = qtb.Filter(k, v)
	}

	for k, v := range exclude {
		qtb = qtb.Exclude(k, v)
	}

	if total, err = qtb.OrderBy("created_at").All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}

// GetFilterSalesOrders : function to get data from database based on parameters with filtered permission
func GetFilterSalesOrders(rq *orm.RequestQuery) (m []*model.SalesOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesOrder))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.SalesOrder
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidSalesOrder : function to check if id is valid in database
func ValidSalesOrder(id int64) (SalesOrder *model.SalesOrder, e error) {
	SalesOrder = &model.SalesOrder{ID: id}
	e = SalesOrder.Read("ID")

	return
}

// Get Grand Total Charge Where SO Status 1, 7, 8
func GetGrandTotalChargeSO(merchant_id int64) (totalCharge float64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q := "SELECT SUM(so.total_charge) total_charge FROM branch b " +
		"LEFT JOIN sales_order so ON so.branch_id = b.id " +
		"WHERE so.status IN (1,7,8) AND b.merchant_id = ?"

	if e = o.Raw(q, merchant_id).QueryRow(&totalCharge); e != nil {
		return totalCharge, e
	}
	return totalCharge, nil
}
