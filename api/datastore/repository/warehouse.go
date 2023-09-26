// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWarehouse find a single data warehouse using field and value condition.
func GetWarehouse(field string, values ...interface{}) (*model.Warehouse, error) {
	m := new(model.Warehouse)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	o.Raw("select * from warehouse_coverage wc where wc.sub_district_id = ? and wc.warehouse_id = ?", m.SubDistrict.ID, m.ID).QueryRows(&m.WarehouseCoverage)
	return m, nil
}

// GetWarehouses get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetWarehouses(rq *orm.RequestQuery) (m []*model.Warehouse, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Warehouse))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Warehouse
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetFilterWarehouse get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterWarehouse(rq *orm.RequestQuery, areaId, warehouseId string) (m []*model.Warehouse, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Warehouse))

	cond := q.GetCond()
	var areaArrDecrypt []int64
	if areaId != "" {
		areaArr := strings.Split(areaId, ",")
		for _, v := range areaArr {
			areaID, _ := common.Decrypt(v)
			areaArrDecrypt = append(areaArrDecrypt, areaID)
		}
		cond1 := orm.NewCondition()
		cond1 = cond1.And("area_id__in", areaArrDecrypt)

		cond = cond.AndCond(cond1)
	}

	var warehouseArrDecrypt []int64
	if warehouseId != "" {
		warehouseArr := strings.Split(warehouseId, ",")
		for _, v := range warehouseArr {
			areaID, _ := common.Decrypt(v)
			warehouseArrDecrypt = append(warehouseArrDecrypt, areaID)
		}
		cond2 := orm.NewCondition()
		cond2 = cond2.And("id__in", warehouseArrDecrypt)

		cond = cond.AndCond(cond2)
	}

	q = q.SetCond(cond)

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Warehouse
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ValidWarehouse : function to check if id is valid in database
func ValidWarehouse(id int64) (w *model.Warehouse, e error) {
	w = &model.Warehouse{ID: id}
	e = w.Read("ID")

	return
}

// CheclWarehousesData : function to get all warehouse data based on filter and exclude parameters
func CheckWarehousesData(filter, exclude map[string]interface{}) (m []*model.Warehouse, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Warehouse))

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

// GetWarehouseSelfPickupByAreaID :
func GetWarehouseSelfPickupByAreaID(areaID int64) (warehouseSelfPickUp *model.Warehouse, e error) {
	var (
		attributeAppConfig       string
		availWarehouseSelfPickUp *model.ConfigApp
	)

	attributeAppConfig = "warehouse_available_self_pickup_"
	// filter attribute base on area
	if areaID == 2 {
		attributeAppConfig += "jkt"
	} else if areaID == 3 {
		attributeAppConfig += "bdg"
	} else if areaID == 4 {
		attributeAppConfig += "smr"
	} else if areaID == 5 {
		attributeAppConfig += "sby"
	}

	if availWarehouseSelfPickUp, e = GetConfigApp("attribute", attributeAppConfig); e != nil {
		return warehouseSelfPickUp, e
	}

	warehouseSelfPickUpID, _ := strconv.ParseInt(availWarehouseSelfPickUp.Value, 10, 64)
	warehouseSelfPickUp = &model.Warehouse{ID: warehouseSelfPickUpID}
	if e = warehouseSelfPickUp.Read("ID"); e != nil {
		return warehouseSelfPickUp, e
	}

	return warehouseSelfPickUp, e
}
