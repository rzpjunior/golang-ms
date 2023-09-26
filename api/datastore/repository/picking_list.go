// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPickingList find a single data picking list using field and value condition.
func GetPickingList(field string, values ...interface{}) (*model.PickingList, error) {
	m := new(model.PickingList)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPickingLists : function to get data from database based on parameters
func GetPickingLists(rq *orm.RequestQuery) (m []*model.PickingList, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")
	q, _ := rq.QueryReadOnly(new(model.PickingList))

	var mx []*model.PickingList
	if _, err = q.Exclude("status__in", 3, 2).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	var resp []*model.PickingList
	for _, v := range mx {
		o.Raw("SELECT so.id, so.branch_id, so.term_payment_sls_id, so.term_invoice_sls_id, so.salesperson_id, so.sales_group_id, so.sub_district_id, so.warehouse_id, so.wrt_id, so.area_id, so.voucher_id, so.price_set_id, so.payment_group_sls_id, so.archetype_id, so.order_type_sls_id, so.order_channel, so.code, so.status, so.recognition_date, so.delivery_date, so.billing_address, so.shipping_address, so.shipping_address_note, so.delivery_fee, so.vou_redeem_code, so.vou_disc_amount, so.point_redeem_amount, so.point_redeem_id, so.total_price, so.total_charge, so.total_weight, so.note, so.reload_packing, so.payment_reminder, so.is_locked, so.has_ext_invoice, so.has_picking_assigned, so.cancel_type, so.created_at, so.created_by, so.last_updated_at, so.last_updated_by, so.finished_at, so.locked_by "+
			"from sales_order so "+
			"join picking_order_assign poa on poa.sales_order_id = so.id "+
			"where poa.picking_list_id = ?", v.ID).QueryRows(&v.SalesOrder)

		countSO := 0
		for _, v2 := range v.SalesOrder {
			if v2.Status == 3 || v2.Status == 4 {
				countSO++
			}

			o.Raw("select count(id) from sales_order_item soi where soi.sales_order_id = ?", v2.ID).QueryRow(&v2.TotalItem)
			v.TotalWeightPickingList += v2.TotalWeight
			v.TotalItemPickingList += v2.TotalItem
			v2.Wrt.Read("ID")
			v2.Branch.Read("ID")
			v2.Branch.Merchant.Read("ID")
			strSepComma := strings.Split(v2.Branch.Merchant.TagCustomer, ",")
			v2.Branch.Merchant.TagCustomerName, _ = util.GetCustomerTag(strSepComma)
		}

		var pickers string
		o.Raw("SELECT sub_picker_id FROM picking_order_assign where picking_list_id = ? limit 1", v.ID).QueryRow(&pickers)

		arrPicker := strings.Split(pickers, ",")
		for _, v2 := range arrPicker {
			pickerID, _ := strconv.Atoi(v2)
			picker := &model.Staff{ID: int64(pickerID)}
			picker.Read("id")

			v.Pickers = append(v.Pickers, picker)
		}

		// this conditions is using if amount of sales order which is not cancelled keep showed
		if countSO != len(v.SalesOrder) {
			resp = append(resp, v)
		}
		v.TotalSalesOrder = len(v.SalesOrder)
		if len(v.SalesOrder) != 0 {
			o.Raw("select s.name from picking_order_assign poa join staff s on s.id = poa.staff_id where poa.sales_order_id = ?", v.SalesOrder[0].ID).QueryRow(&v.PickerName)
		}
	}
	total = int64(len(resp))
	return resp, total, nil
}

// ValidPickingList : function to check if id is valid in database
func ValidPickingList(id int64) (pickingOrder *model.PickingList, e error) {
	pickingOrder = &model.PickingList{ID: id}
	e = pickingOrder.Read("ID")

	return
}

// GetSalesOrderbyPickingList : function to get all SO based on a picking list
func GetSalesOrderbyPickingList(rq *orm.RequestQuery) (m []*model.SalesOrder, pl *model.PickingList, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	mapSO := map[int64]bool{}
	q, _ := rq.QueryReadOnly(new(model.PickingOrderAssign))

	var mx []*model.PickingOrderAssign
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, nil, err
	}

	pl = &model.PickingList{ID: mx[0].PickingList.ID}
	if err = pl.Read("id"); err != nil {
		return nil, nil, err
	}

	for _, v := range mx {
		if !mapSO[v.SalesOrder.ID] {
			mapSO[v.SalesOrder.ID] = true
			v.SalesOrder.Read("ID")
			v.SalesOrder.Branch.Read("ID")
			v.SalesOrder.Branch.Merchant.Read("ID")
			v.SalesOrder.Wrt.Read("ID")
			m = append(m, v.SalesOrder)
		}
	}

	return m, pl, nil
}
