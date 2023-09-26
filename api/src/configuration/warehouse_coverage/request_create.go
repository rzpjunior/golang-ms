// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse_coverage

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold warehouse coverage set request data
type createRequest struct {
	SubDistrictID string `json:"sub_district_id" valid:"required"`
	WarehouseID   string `json:"warehouse_id" valid:"required"`
	MainWarehouse bool   `json:"main_warehouse"`

	MainWarehouseInt int8               `json:"-"`
	SubDistrict      *model.SubDistrict `json:"-"`
	Warehouse        *model.Warehouse   `json:"-"`
	WarehouseType    *model.Glossary    `json:"-"`

	WarehouseCoverage         *model.WarehouseCoverage `json:"-"`
	AffectedWarehouseCoverage *model.WarehouseCoverage `json:"-"`
	ChangeMainWarehouse       bool                     `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate warehouse coverage request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var filter, exclude map[string]interface{}

	subDistrictId, e := common.Decrypt(c.SubDistrictID)
	if e != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	}
	if c.SubDistrict, e = repository.ValidSubDistrict(subDistrictId); e != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	}

	// warehouse validation
	whID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}
	if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	// check if the warehouse is already inside the subdistrict
	filter = map[string]interface{}{"sub_district_id": c.SubDistrict.ID, "warehouse_id": c.Warehouse.ID}
	if _, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse != 0 {
		o.Failure("warehouse_coverage.invalid", util.ErrorExistWarehouseCoverage())
	}

	if c.SubDistrict.Area.ID != c.Warehouse.Area.ID {
		o.Failure("area_id.invalid", util.ErrorMustBeSame("warehouse area", "sub district area"))
	}

	if c.WarehouseType, e = repository.GetGlossaryMultipleValue("table", "warehouse", "attribute", "warehouse_type", "value_int", c.Warehouse.WarehouseType); e != nil {
		o.Failure("warehouse_type.invalid", util.ErrorInvalidData("warehouse type"))
		return o
	}
	if c.WarehouseType.ValueName == "HUB" && c.MainWarehouse == true {
		o.Failure("main_warehouse.invalid", util.HubMainWarehouse())
	} else if c.WarehouseType.ValueName == "HUB" && c.MainWarehouse == false {
		// check if the parent warehouse serving the sub district
		filter = map[string]interface{}{"sub_district_id": c.SubDistrict.ID, "warehouse_id": c.Warehouse.ParentID}
		if _, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse == 0 {
			o.Failure("parent_warehouse.invalid", util.HubNeedParent())
		}

		// check if there's another hub in that sub district with same parent
		filter = map[string]interface{}{"sub_district_id": c.SubDistrict.ID, "parent_warehouse_id": c.Warehouse.ParentID}
		if _, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse != 0 {
			o.Failure("warehouse_coverage.invalid", util.HubOnlyOne())
		}
	} else if c.WarehouseType.ValueName != "HUB" && c.MainWarehouse == true {
		// check if there's a main warehouse in the sub district
		filter = map[string]interface{}{"sub_district_id": c.SubDistrict.ID, "main_warehouse": 1}
		if data, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err == nil && countWarehouse != 0 {
			c.ChangeMainWarehouse = true
			c.AffectedWarehouseCoverage = data[0]
		}
	}

	if c.MainWarehouse {
		c.MainWarehouseInt = 1
	} else {
		c.MainWarehouseInt = 2
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"sub_district_id.required": util.ErrorInputRequired("sub district id"),
		"warehouse_id.required":    util.ErrorInputRequired("warehouse id"),
	}
}
