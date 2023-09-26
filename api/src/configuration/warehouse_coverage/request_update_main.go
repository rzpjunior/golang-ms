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

// updateMainRequest : struct to hold warehouse coverage set request data
type updateMainRequest struct {
	ID int64 `json:"-"`

	WarehouseCoverage         *model.WarehouseCoverage `json:"-"`
	AffectedWarehouseCoverage *model.WarehouseCoverage `json:"-"`
	ChangeMainWarehouse       bool                     `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate warehouse coverage request data
func (c *updateMainRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var filter, exclude map[string]interface{}
	var e error

	if c.WarehouseCoverage, e = repository.ValidWarehouseCoverage(c.ID); e != nil {
		o.Failure("warehouse_coverage_id.invalid", util.ErrorInvalidData("warehouse coverage"))
	}

	// cannot delete if warehouse is a main warehouse
	if c.WarehouseCoverage.MainWarehouse == 1 {
		o.Failure("warehouse_coverage.invalid", util.IsMainWarehouse())
	}

	// check if the warehouse is a hub
	if e = c.WarehouseCoverage.Warehouse.Read("ID"); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		return o
	}
	warehouseType, e := repository.GetGlossaryMultipleValue("table", "warehouse", "attribute", "warehouse_type", "value_int", c.WarehouseCoverage.Warehouse.WarehouseType)
	if e != nil {
		o.Failure("warehouse_type.invalid", util.ErrorInvalidData("warehouse type"))
		return o
	}
	if warehouseType.ValueName == "HUB" {
		o.Failure("main_warehouse.invalid", util.HubMainWarehouse())
	}

	// check if the subdistrict has another main warehouse
	filter = map[string]interface{}{"sub_district_id": c.WarehouseCoverage.SubDistrict, "main_warehouse": 1}
	if data, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse != 0 {
		c.ChangeMainWarehouse = true
		c.AffectedWarehouseCoverage = data[0]
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateMainRequest) Messages() map[string]string {
	return map[string]string{}
}
