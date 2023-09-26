// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// updateShippingAddressRequest : struct to hold price set request data
type updateShippingAddressRequest struct {
	ID int64 `json:"-" valid:"required"`
	//AddressName             string `json:"address_name" valid:"required"`
	RecipientName           string `json:"recipient_name" valid:"required"`
	RecipientPhoneNumber    string `json:"recipient_phone_number" valid:"required"`
	RecipientAltPhoneNumber string `json:"recipient_alt_phone_number"`
	ShippingAddress         string `json:"shipping_address" valid:"required"`
	ShippingNote            string `json:"shipping_note"`
	AreaId                  string `json:"area_id" valid:"required"`
	SubDistrictId           string `json:"sub_district_id" valid:"required"`
	//WarehouseId             string `json:"warehouse_id" valid:"required"`

	Area              *model.Area
	SubDistrict       *model.SubDistrict
	WarehouseCoverage *model.WarehouseCoverage

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *updateShippingAddressRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	//var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if areaId, err := common.Decrypt(c.AreaId); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, err = repository.ValidArea(areaId); err != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area_id.active", util.ErrorActive("area"))
			}
		}
	}

	if subDistrictId, err := common.Decrypt(c.SubDistrictId); err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	} else {
		if c.SubDistrict, err = repository.ValidSubDistrict(subDistrictId); err != nil {
			o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
		} else {
			if c.SubDistrict.Status != int8(1) {
				o.Failure("sub_district_id.active", util.ErrorActive("sub district"))
			}
		}
	}

	if c.SubDistrict != nil {
		orSelect.Raw("select * from warehouse_coverage wc where sub_district_id =? and main_warehouse = 1", c.SubDistrict.ID).QueryRow(&c.WarehouseCoverage)
	}

	if c.WarehouseCoverage == nil {
		o.Failure("warehouse_coverage.invalid", util.ErrorMustBeSame("warehouse_coverage", "warehouse"))
	}

	if len(c.ShippingNote) > 100 {
		o.Failure("shipping_note.invalid", util.ErrorCharLength("shipping note", 100))
	}

	if len(c.ShippingAddress) > 350 {
		o.Failure("shipping_address.invalid", util.ErrorCharLength("shipping address", 350))
	}

	//if warehouseId, err := common.Decrypt(c.WarehouseId); err != nil {
	//	o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	//} else {
	//	if c.DefaultWarehouse, err = repository.ValidWarehouse(warehouseId); err != nil {
	//		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	//	} else {
	//		if c.DefaultWarehouse.Status != int8(1) {
	//			o.Failure("warehouse_id.active", util.ErrorActive("warehouse"))
	//		}
	//	}
	//}

	return o
}

// Messages : function to return error validation messages
func (c *updateShippingAddressRequest) Messages() map[string]string {
	return map[string]string{
		"recipient_name.required":         util.ErrorInputRequired("recipient name"),
		"recipient_phone_number.required": util.ErrorInputRequired("recipient phone number"),
		"shipping_address.required":       util.ErrorInputRequired("shipping address"),
		"area_id.required":                util.ErrorInputRequired("area"),
		"sub_district_id.required":        util.ErrorInputRequired("sub district id"),
		//"warehouse_id.required":           util.ErrorInputRequired("default warehouse id"),
	}
}
