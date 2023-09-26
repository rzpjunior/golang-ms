// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDeliveryReturn find a single data price set using field and value condition.
func GetDeliveryReturn(field string, values ...interface{}) (*model.DeliveryReturn, error) {
	m := new(model.DeliveryReturn)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "DeliveryReturnItems", 2)

	for _, row := range m.DeliveryReturnItems {
		o.Raw("select value_name from glossary where `table` = ? and `attribute` = ? and `value_int` = ?", "all", "waste_reason", row.WasteReason).QueryRow(&row.WasteReasonValue)
	}

	return m, nil
}

// GetDeliveryReturns : function to get data from database based on parameters
func GetDeliveryReturns(rq *orm.RequestQuery) (m []*model.DeliveryReturn, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryReturn))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.DeliveryReturn
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterDeliveryReturns : function to get data from database based on parameters with filtered permission
func GetFilterDeliveryReturns(rq *orm.RequestQuery) (m []*model.DeliveryReturn, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.DeliveryReturn))

	if total, err = q.Filter("status__in", 1, 2, 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.DeliveryReturn
	if _, err = q.Filter("status__in", 1, 2, 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidDeliveryReturn : function to check if id is valid in database
func ValidDeliveryReturn(id int64) (deliveryReturn *model.DeliveryReturn, e error) {
	deliveryReturn = &model.DeliveryReturn{ID: id}
	e = deliveryReturn.Read("ID")

	return
}

// CheckDeliveryReturnProductStatus : function to check if product is valid in table
func CheckDeliveryReturnProductStatus(productID int64, status int8, warehouseArr ...string) (*model.DeliveryReturnItem, int64, error) {
	var err error
	o := orm.NewOrm()
	dri := new(model.DeliveryReturnItem)
	o.Using("read_only")
	q := o.QueryTable(dri).RelatedSel("DeliveryReturn").Filter("DeliveryReturn__Status", status).Filter("product_id", productID)

	if len(warehouseArr) > 0 {
		q = q.Exclude("DeliveryReturn__Warehouse__id__in", warehouseArr)
	}

	if total, err := q.All(dri); err == nil {
		return dri, total, nil
	}

	return nil, 0, err
}
