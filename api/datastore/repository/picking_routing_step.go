// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPickingRoutingStep : find a single data using field and value condition.
func GetPickingRoutingStep(field string, values ...interface{}) (*model.PickingRoutingStep, error) {
	m := new(model.PickingRoutingStep)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetPickingRoutingSteps : function to get data from database based on parameters
func GetPickingRoutingSteps(rq *orm.RequestQuery) (m []*model.PickingRoutingStep, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PickingRoutingStep))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PickingRoutingStep
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			_, err = o.Raw("SELECT * from sales_order so "+
				"join picking_order_assign poa on poa.sales_order_id = so.id "+
				"where poa.picking_list_id = ?", v.PickingList).QueryRows(&v.PickingList.SalesOrder)
			if err != nil {
				return nil, total, err
			}

			for _, v2 := range v.PickingList.SalesOrder {
				if err = v2.Branch.Read("id"); err != nil {
					return nil, total, err
				}
				if err = v2.Branch.Merchant.Read("id"); err != nil {
					return nil, total, err
				}
			}
		}

		if err = o.Raw("select * from staff s join picking_order_assign poa on poa.staff_id=s.id where poa.picking_list_id = ? LIMIT 1", mx[0].PickingList).QueryRow(&mx[0].LeadPicker); err != nil {
			return nil, total, err
		}
		return mx, total, nil
	}

	return nil, total, err
}

func CheckPickingRoutingStepData(filter, exclude map[string]interface{}) (steps []*model.PickingRoutingStep, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PickingRoutingStep))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&steps); err == nil {
		return steps, total, nil
	}

	return nil, 0, err
}

func ValidPickingRoutingStep(id int64) (steps *model.PickingRoutingStep, e error) {
	pickingRoutingStep := &model.PickingRoutingStep{ID: id}
	e = pickingRoutingStep.Read("ID")

	return
}

// GetFirstPickingRoutingStep : function to get data of picking routing step that has status in progress and with the lowest sequence number
func GetFirstPickingRoutingStep(staffId, pickingListId int64) (m *model.PickingRoutingStep, e error) {
	o := orm.NewOrm()
	o.Using("read_only")
	mdb := mongodb.NewMongo()

	if e = o.Raw("SELECT * FROM picking_routing_step WHERE staff_id = ? and picking_list_id = ? and status_step = 2 ORDER BY sequence ASC LIMIT 1", staffId, pickingListId).QueryRow(&m); e != nil {
		return nil, e
	}

	o.LoadRelated(m, "bin", 0)
	o.LoadRelated(m, "pickingorderitem", 0)
	o.LoadRelated(m, "pickinglist", 0)

	if e = o.Raw("select * from staff s join picking_order_assign poa on poa.staff_id=s.id where poa.picking_list_id = ?", m.PickingList).QueryRow(&m.LeadPicker); e != nil {
		return nil, e
	}

	if m.StepType == 2 || m.StepType == 3 {
		o.LoadRelated(m.PickingOrderItem, "pickingorderassign", 0)
		o.LoadRelated(m.PickingOrderItem, "product", 0)
		o.LoadRelated(m.PickingOrderItem.Product, "uom", 0)
		o.LoadRelated(m.PickingOrderItem.Product, "ProductImage", 0)
		o.LoadRelated(m.PickingOrderItem.PickingOrderAssign, "salesorder", 0)
		o.LoadRelated(m.PickingOrderItem.PickingOrderAssign.SalesOrder, "branch", 0)
		o.LoadRelated(m.PickingOrderItem.PickingOrderAssign.SalesOrder.Branch, "merchant", 0)

		strSepComma := strings.Split(m.PickingOrderItem.PickingOrderAssign.SalesOrder.Branch.Merchant.TagCustomer, ",")
		m.PickingOrderItem.PickingOrderAssign.SalesOrder.Branch.Merchant.TagCustomerName, _ = util.GetCustomerTag(strSepComma)
		o.Raw("SELECT note FROM sales_order_item WHERE sales_order_id=? AND product_id=?", m.PickingOrderItem.PickingOrderAssign.SalesOrder.ID, m.PickingOrderItem.Product.ID).QueryRow(&m.SalesOrderItemNote)

		filter := bson.D{
			{"sales_order_id", m.PickingOrderItem.PickingOrderAssign.SalesOrder.ID},
			{"product_id", m.PickingOrderItem.Product.ID},
			{"pack_type", bson.M{"$ne": -1}},
			{"status", 1},
		}

		var res []byte
		if res, e = mdb.GetAllDataWithFilter("Packing_Sales_Order", filter); e != nil {
			mdb.DisconnectMongoClient()
			return nil, e
		}

		if len(res) == 0 {
			mdb.DisconnectMongoClient()
			return m, nil
		}
		var pr []*model.PackRecommendation
		// region convert byte data to json data
		if e = json.Unmarshal(res, &pr); e != nil {
			mdb.DisconnectMongoClient()
			return nil, e
		}
		// endregion
		m.PackRecommendation = pr

		for _, v := range pr {
			pcIDStr := common.Encrypt(v.PackingOrderID)
			pcID, _ := strconv.Atoi(pcIDStr)
			v.PackingOrderID = int64(pcID)

			pIDStr := common.Encrypt(v.ProductID)
			pID, _ := strconv.Atoi(pIDStr)
			v.ProductID = int64(pID)

			soIDStr := common.Encrypt(v.SalesOrderID)
			soID, _ := strconv.Atoi(soIDStr)
			v.SalesOrderID = int64(soID)
		}
	}

	mdb.DisconnectMongoClient()
	return m, nil
}
