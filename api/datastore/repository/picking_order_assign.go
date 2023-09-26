// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPickingOrderAssign find a single data picking order assign using field and value condition.
func GetPickingOrderAssign(field string, values ...interface{}) (*model.PickingOrderAssign, error) {
	m := new(model.PickingOrderAssign)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}

	o.Raw("select count(*) from sales_order_item soi where soi.sales_order_id = ?", m.SalesOrder.ID).QueryRow(&m.TotalItemSO)
	o.Raw("select count(*) from picking_order_item poi where poi.picking_order_assign_id = ? and poi.pick_qty > 0", m.ID).QueryRow(&m.TotalItemOnProgress)
	o.Raw("select si.delta_print from sales_invoice si where si.sales_order_id  = ?", m.SalesOrder.ID).QueryRow(&m.DeltaPrintSalesInvoice)
	o.Raw("select do.delta_print from delivery_order do where do.sales_order_id  = ?", m.SalesOrder.ID).QueryRow(&m.DeltaPrintDeliveryOrder)

	o.LoadRelated(m, "PickingOrderItem", 2)
	m.SalesOrder.Branch.Merchant.Read("ID")
	for _, p := range m.PickingOrderItem {
		o.Raw("select note from sales_order_item where product_id = ? and sales_order_id = ?", p.Product.ID, m.SalesOrder.ID).QueryRow(&p.SalesOrderItemNote)
		o.LoadRelated(p.Product, "ProductImage", 1)
	}

	return m, nil
}

// GetPickingOrderAssigns : function to get data from database based on parameters
func GetPickingOrderAssigns(rq *orm.RequestQuery) (m []*model.PickingOrderAssign, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.PickingOrderAssign))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PickingOrderAssign
	if _, err = q.RelatedSel(2).All(&mx, rq.Fields...); err == nil {
		for _, pa := range mx {
			pa.SalesOrder.Branch.Merchant.Read("ID")
			o.Raw("select count(soi.id) from sales_order_item soi where soi.sales_order_id = ?", pa.SalesOrder.ID).QueryRow(&pa.TotalItemSO)
			o.Raw("select count(poi.id) from picking_order_item poi where poi.picking_order_assign_id = ? and poi.pick_qty > 0", pa.ID).QueryRow(&pa.TotalItemOnProgress)
			o.Raw("select si.delta_print from sales_invoice si where si.sales_order_id  = ?", pa.SalesOrder.ID).QueryRow(&pa.DeltaPrintSalesInvoice)
			o.Raw("select do.delta_print from delivery_order do where do.sales_order_id  = ?", pa.SalesOrder.ID).QueryRow(&pa.DeltaPrintDeliveryOrder)

			strSepComma := strings.Split(pa.SalesOrder.Branch.Merchant.TagCustomer, ",")
			pa.SalesOrder.Branch.Merchant.TagCustomerName, _ = util.GetCustomerTag(strSepComma)
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetPickingOrderAssignsforDispatch : function to get data from database based on parameters
func GetPickingOrderAssignsforDispatch(rq *orm.RequestQuery) (m []*model.PickingOrderAssign, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.PickingOrderAssign))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PickingOrderAssign
	if _, err = q.RelatedSel(2).All(&mx, rq.Fields...); err == nil {
		for _, pa := range mx {
			pa.SalesOrder.Branch.Merchant.Read("ID")
		}
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPickingOrderAssign : function to check if id is valid in database
func ValidPickingOrderAssign(id int64) (pickingOrderAssign *model.PickingOrderAssign, e error) {
	pickingOrderAssign = &model.PickingOrderAssign{ID: id}
	e = pickingOrderAssign.Read("ID")

	return
}

// ValidPickingOrderItem : function to check if id is valid in database
func ValidPickingOrderItem(id int64) (pickingOrderItem *model.PickingOrderItem, e error) {
	pickingOrderItem = &model.PickingOrderItem{ID: id}
	e = pickingOrderItem.Read("ID")

	return
}

// CheckPickingOrderAssignData : function to check PickingOrderAssign data based on filter and exclude parameters
func CheckPickingOrderAssignData(filter, exclude map[string]interface{}) (PickingOrderAssign []*model.PickingOrderAssign, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PickingOrderAssign))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&PickingOrderAssign); err == nil {
		return PickingOrderAssign, total, nil
	}

	return nil, 0, err
}

// GetItemByPickingAssignId : get all item data based on picking assign id
func GetItemByPickingAssignId(field string, values ...interface{}) (pis []*model.PickingOrderItem, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	m := new(model.PickingOrderAssign)

	if err = o.QueryTable(m).Filter(field, values...).Limit(1).One(m); err != nil {
		return nil, err
	}

	if _, err = o.Raw("SELECT * FROM picking_order_item WHERE picking_order_assign_id = ?", m.ID).QueryRows(&pis); err != nil {
		return nil, err
	}

	for _, v := range pis {
		v.Product.Read("ID")
		v.Product.Uom.Read("ID")
	}

	return pis, err
}
