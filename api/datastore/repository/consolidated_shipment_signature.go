// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// ValidConsolidatedShipmentSignature : function to check if id is valid in database
func ValidConsolidatedShipmentSignature(id int64) (consolidatedShipmentSignature *model.ConsolidatedShipmentSignature, e error) {
	consolidatedShipmentSignature = &model.ConsolidatedShipmentSignature{ID: id}
	e = consolidatedShipmentSignature.Read("ID")

	return
}

// IsRoleCSAlreadySigned : check if role already sign consolidated shipment
func IsRoleCSAlreadySigned(consolidatedShipmentID int64, jobFunction string) (bool, error) {
	consolidatedShipmentSignature := new(model.ConsolidatedShipmentSignature)
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err := o.QueryTable(consolidatedShipmentSignature).Filter("consolidated_shipment_id", consolidatedShipmentID).Filter("job_function", jobFunction).Count(); err != nil || total != 0 {
		return false, err
	}

	return true, nil
}

// GetFilterConsolidatedShipmentSignature : function to check data based on filter and exclude parameters
func GetFilterConsolidatedShipmentSignature(filter, exclude map[string]interface{}) (m []*model.ConsolidatedShipmentSignature, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.ConsolidatedShipmentSignature))
	orm := orm.NewOrm()
	orm.Using("read_only")

	for i, v := range filter {
		o = o.Filter(i, v)
	}

	for i, v := range exclude {
		o = o.Exclude(i, v)
	}

	total, err = o.All(&m)
	if err != nil {
		return nil, 0, err
	}

	return m, total, nil
}
