// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPickingOrder find a single data picking order using field and value condition.
func GetPickingOrder(field string, values ...interface{}) (*model.PickingOrder, error) {
	m := new(model.PickingOrder)
	o := orm.NewOrm()
	o.Using("read_only")
	var filter, exclude map[string]interface{}

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "PickingOrderAssign", 2)
	for _, pa := range m.PickingOrderAssign {
		pa.SalesOrder.Read("ID")
		pa.SalesOrder.SubDistrict.Read("ID")
		pa.SalesOrder.SubDistrict.District.Read("ID")
		pa.SalesOrder.Branch.Read("ID")
		pa.SalesOrder.Branch.Archetype.Read("ID")
		pa.SalesOrder.Branch.Archetype.BusinessType.Read("ID")
		o.Raw("select count(*) from sales_order_item soi where soi.sales_order_id = ?", pa.SalesOrder.ID).QueryRow(&pa.TotalItemSO)
		o.Raw("select count(*) from picking_order_item poi where poi.picking_order_assign_id = ? and poi.pick_qty > 0", pa.ID).QueryRow(&pa.TotalItemOnProgress)

		filter = map[string]interface{}{"picking_list_id": pa.PickingList.ID, "status_step__in": []int64{2, 3}}
		_, total, err := CheckPickingRoutingStepData(filter, exclude)
		if err != nil {
			return m, nil
		}

		if total == 0 {
			pa.PickingList.PickingRouting = 3
		} else {
			pa.PickingList.PickingRouting = 1
		}

	}

	return m, nil
}

// GetPickingOrders : function to get data from database based on parameters
func GetPickingOrders(rq *orm.RequestQuery) (m []*model.PickingOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PickingOrder))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PickingOrder
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPickingOrder : function to check if id is valid in database
func ValidPickingOrder(id int64) (pickingOrder *model.PickingOrder, e error) {
	pickingOrder = &model.PickingOrder{ID: id}
	e = pickingOrder.Read("ID")

	return
}

// CheckPickingOrderData : function to check PickingOrder data based on filter and exclude parameters
func CheckPickingOrderData(filter, exclude map[string]interface{}) (PickingOrder []*model.PickingOrder, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PickingOrder))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&PickingOrder); err == nil {
		return PickingOrder, total, nil
	}

	return nil, 0, err
}
