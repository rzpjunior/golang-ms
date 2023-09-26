// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// ValidPurchasePlan : function to check if id is valid in database
func ValidPurchasePlan(id int64) (purchasePlan *model.PurchasePlan, e error) {
	purchasePlan = &model.PurchasePlan{ID: id}
	e = purchasePlan.Read("ID")

	return
}

// GetPurchasePlans : function to get list data of purchase plan & load purchase plan item
func GetPurchasePlans(rq *orm.RequestQuery, warehouseID int64) (mx []*model.PurchasePlan, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	q, _ := rq.QueryReadOnly(new(model.PurchasePlan))

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)

	if warehouseID == warehouse.ID || warehouseID == 0 {
		if total, err = q.Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if total, err = q.Filter("warehouse_id", warehouseID).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).Filter("warehouse_id", warehouseID).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	}

	for _, item := range mx {

		err = o.Raw("SELECT count(id) FROM purchase_plan_item WHERE purchase_plan_id = ?", item.ID).QueryRow(&item.TotalSku)

		if err != nil {
			return nil, total, err
		}

		_, err = o.Raw("SELECT u.name as uom_name, sum(ppi.weight) as total_weight FROM purchase_plan_item ppi JOIN product p ON p.id =  ppi.product_id JOIN uom u ON u.id = p.uom_id WHERE purchase_plan_id = ? GROUP BY p.uom_id", item.ID).QueryRows(&item.TonnagePurchasePlan)

		if err != nil {
			return nil, total, err
		}
	}

	return mx, total, nil
}

// GetPurchasePlan find a single data price set using field and value condition.
func GetPurchasePlan(field string, values ...interface{}) (*model.PurchasePlan, error) {
	m := new(model.PurchasePlan)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "PurchasePlanItems", 2)

	for _, v := range m.PurchasePlanItems {
		o.LoadRelated(v, "PurchaseOrderItems", 2)
	}

	return m, nil
}

// GetPurchasePlansInFieldPurchaserApps : function to get list data of purchase plan & load purchase plan item
func GetPurchasePlansInFieldPurchaserApps(rq *orm.RequestQuery, warehouseID int64) (mx []*model.PurchasePlan, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)
	q, _ := rq.QueryReadOnly(new(model.PurchasePlan))

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)

	if warehouseID == warehouse.ID || warehouseID == 0 {
		if total, err = q.Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if total, err = q.Filter("warehouse_id", warehouseID).Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.RelatedSel(2).Filter("warehouse_id", warehouseID).All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	}

	for _, item := range mx {

		err = o.Raw("SELECT count(id) FROM purchase_plan_item WHERE purchase_plan_id = ?", item.ID).QueryRow(&item.TotalSku)

		if err != nil {
			return nil, total, err
		}

		_, err = o.Raw("SELECT u.name as uom_name, sum(ppi.purchase_plan_qty) as total_weight FROM purchase_plan_item ppi JOIN product p ON p.id =  ppi.product_id JOIN uom u ON u.id = p.uom_id WHERE purchase_plan_id = ? GROUP BY p.uom_id", item.ID).QueryRows(&item.TonnagePurchasePlan)

		if err != nil {
			return nil, total, err
		}
	}

	return mx, total, nil
}

// CheckPurchasePlanData : function to get all purchase plan data based on filter and exclude parameters
func CheckPurchasePlanData(filter, exclude map[string]interface{}) (m []*model.PurchasePlan, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PurchasePlan))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}

// GetFilterPurchasePlans : function to get data from database based on parameters with filtered permission
func GetFilterPurchasePlans(rq *orm.RequestQuery) (mx []*model.PurchasePlan, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PurchasePlan))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	return mx, total, nil
}
