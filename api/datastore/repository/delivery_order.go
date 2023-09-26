// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDeliveryOrder find a single data price set using field and value condition.
func GetDeliveryOrder(field string, values ...interface{}) (*model.DeliveryOrder, error) {
	m := new(model.DeliveryOrder)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "DeliveryOrderItems", 2)

	o.Raw("SELECT * FROM sales_invoice where sales_order_id = ?", m.SalesOrder.ID).QueryRows(&m.SalesInvoice)
	o.Raw("select note from audit_log where `type` = 'delivery_order' and `function` = 'cancel' and ref_id = ?", m.ID).QueryRow(&m.CancellationNote)

	return m, nil
}

// GetDeliveryOrders : function to get data from database based on parameters
func GetDeliveryOrders(rq *orm.RequestQuery) (m []*model.DeliveryOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryOrder))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.DeliveryOrder
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterDeliveryOrders : function to get data from database based on parameters with filtered permission
func GetFilterDeliveryOrders(rq *orm.RequestQuery) (m []*model.DeliveryOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryOrder))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.DeliveryOrder
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidDeliveryOrder : function to check if id is valid in database
func ValidDeliveryOrder(id int64) (deliveryOrder *model.DeliveryOrder, e error) {
	deliveryOrder = &model.DeliveryOrder{ID: id}
	e = deliveryOrder.Read("ID")

	return
}

// CheckDeliveryOrderProductStatus : function to check if product is valid in table
func CheckDeliveryOrderProductStatus(productID int64, status []int8, warehouseArr ...string) (*model.DeliveryOrderItem, int64, error) {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	doi := new(model.DeliveryOrderItem)
	q := o.QueryTable(doi).RelatedSel("DeliveryOrder").Filter("DeliveryOrder__Status__in", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("DeliveryOrder__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(doi); err == nil {
		return doi, total, nil
	}

	return nil, 0, err
}

// GetDeliveryOrderItems : function to get data from database based on parameters
func GetDeliveryOrderItems(rq *orm.RequestQuery) (gris []*model.DeliveryOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryOrderItem))

	if total, err = q.All(&gris, rq.Fields...); err == nil {
		return gris, total, nil
	}

	return nil, total, err
}

// GetFilterDeliveryOrderItems : function to get data from database based on parameters with filtered permission
func GetFilterDeliveryOrderItems(rq *orm.RequestQuery) (gris []*model.DeliveryOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryOrderItem))

	if total, err = q.All(&gris, rq.Fields...); err == nil {
		return gris, total, nil
	}

	return nil, total, err
}

// GetDeliveryOrdersForPrint : function to get data from database based on parameters
func GetDeliveryOrdersForPrint(rq *orm.RequestQuery) (m []*model.DeliveryOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryOrder))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Exclude("status__in", 4, 2, 3).Count(); err != nil {
		return nil, total, err
	}

	var mx []*model.DeliveryOrder
	if _, err = q.Exclude("status__in", 4, 2, 3).RelatedSel(2).All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {
		o.LoadRelated(v, "DeliveryOrderItems", 2)
		v.SalesOrder.Branch.Merchant.Read("ID")
		o.Raw("SELECT * FROM sales_invoice where sales_order_id = ?", v.SalesOrder.ID).QueryRows(&v.SalesInvoice)
		o.Raw("select note from audit_log where `type` = 'delivery_order' and `function` = 'cancel' and ref_id = ?", v.ID).QueryRow(&v.CancellationNote)

	}
	return mx, total, nil

}
