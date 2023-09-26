// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse_coverage

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// deleteRequest : struct to hold warehouse coverage set request data
type deleteRequest struct {
	ID int64 `json:"-"`

	WarehouseCoverage *model.WarehouseCoverage `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate warehouse coverage request data
func (c *deleteRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var filter, exclude map[string]interface{}
	var e error

	if c.WarehouseCoverage, e = repository.ValidWarehouseCoverage(c.ID); e != nil {
		o.Failure("warehouse_coverage_id.invalid", util.ErrorInvalidData("warehouse coverage"))
	}

	if c.WarehouseCoverage.MainWarehouse == 1 {
		o.Failure("warehouse_coverage.invalid", util.MainWarehouseDelete())
	}

	// check if the warehouse has a hub
	filter = map[string]interface{}{"sub_district_id": c.WarehouseCoverage.SubDistrict, "parent_warehouse_id": c.WarehouseCoverage.Warehouse}
	if _, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse != 0 {
		o.Failure("warehouse_coverage.invalid", util.HubStillExist())
	}

	return o
}

// Messages : function to return error validation messages
func (c *deleteRequest) Messages() map[string]string {
	return map[string]string{}
}
