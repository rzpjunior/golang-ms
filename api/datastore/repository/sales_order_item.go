// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetSalesOrderItem find a single data sales term using field and value condition.
func GetSalesOrderItem(field string, values ...interface{}) (*model.SalesOrderItem, error) {
	m := new(model.SalesOrderItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetSalesOrderItems : function to get data from database based on parameters
func GetSalesOrderItems(rq *orm.RequestQuery, isReport ...int8) (m []*model.SalesOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesOrderItem))

	if total, err = q.All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// GetFilterSalesOrderItems : function to get data from database based on parameters with filtered permission
func GetFilterSalesOrderItems(rq *orm.RequestQuery, isReport ...int8) (m []*model.SalesOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.SalesOrderItem))

	if total, err = q.All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// ValidSalesOrderItem : function to check if id is valid in database
func ValidSalesOrderItem(id int64) (soi *model.SalesOrderItem, e error) {
	soi = &model.SalesOrderItem{ID: id}
	e = soi.Read("ID")

	return
}

func GetSalesOrderItemRecapReport(where map[string]interface{}) (soir []*model.SalesOrderItem, total int64, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	o.Raw("select so.code order_code, p.code product_code, p.name product_name, c.name product_category, uom.name uom, soi.note order_item_note, " +
		"soi.order_qty ordered_qty, sii.invoice_qty invoice_qty, soi.unit_price order_unit_price, soi.shadow_price order_unit_shadow_price, " +
		"soi.subtotal subtotal, " +
		"soi.weight total_weight, so.recognition_date order_date, so.delivery_date order_delivery_date, a.name area, w.name warehouse, " +
		"wrt.name wrt, so.status order_status " +
		"from sales_order_item soi" +
		"join sales_order so ON so.id = soi.sales_order_id" +
		"join sales_invoice_item sii ON sii.sales_order_item_id = soi.id" +
		"join wrt ON wrt.id = so.wrt_id" +
		"join warehouse w ON w.id = so.warehouse_id" +
		"join area a ON a.id = so.area_id" +
		"join product p ON p.id = soi.product_id" +
		"join uom ON uom.id = p.uom_id" +
		"join category c ON c.id = p.category_id").QueryRow(&soir)
	return
}

// CheckSalesOrderItemData : function to check data based on filter and exclude parameters
func CheckSalesOrderItemData(filter, exclude map[string]interface{}) (soi []*model.SalesOrderItem, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.SalesOrderItem))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if countResult, err := o.All(&soi); err == nil {
		return soi, countResult, nil
	}

	return nil, 0, err
}
