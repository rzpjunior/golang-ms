// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sub_district

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetSubDistrict get all data sales_order_item that matched with query request parameters.
// returning slices of SalesOrderDeleted, total data without limit and error.
func GetSubDistrict(rq *orm.RequestQuery) (m []*model.AdmDivision, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.AdmDivision))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.AdmDivision
	if _, err = q.RelatedSel().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
