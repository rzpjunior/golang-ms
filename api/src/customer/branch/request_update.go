// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID                 int64  `json:"-" valid:"required"`
	PriceSetID         string `json:"price_set_id" valid:"required"`
	PicName            string `json:"pic_name" valid:"required"`
	PhoneNumber        string `json:"phone_number" valid:"required"`
	AltPhoneNumber     string `json:"alt_phone_number"`
	ShippingAddress    string `json:"shipping_address" valid:"required"`
	Note               string `json:"note"`
	SubDistrictID      string `json:"sub_district_id" valid:"required"`
	AreaID             string `json:"area_id" valid:"required"`
	WarehouseID        string `json:"warehouse_id" valid:"required"`
	NotePriceSetChange string `json:"-"`

	PriceSet    *model.PriceSet    `json:"-"`
	SubDistrict *model.SubDistrict `json:"-"`
	Area        *model.Area
	Warehouse   *model.Warehouse `json:"-"`
	Branch      *model.Branch    `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var subDistrictID int64

	c.Branch = &model.Branch{ID: c.ID}
	if err = c.Branch.Read("ID"); err != nil {
		o.Failure("branch.invalid", util.ErrorInvalidData("branch"))
	}

	if err = c.Branch.PriceSet.Read("ID"); err != nil {
		o.Failure("price_set.invalid", util.ErrorInvalidData("price set"))
	}

	areaID, _ := common.Decrypt(c.AreaID)
	c.Area = &model.Area{ID: areaID}
	if err = c.Area.Read("ID"); err != nil {
		o.Failure("branch_area_id", util.ErrorInvalidData("branch area"))
	}

	priceSetID, _ := common.Decrypt(c.PriceSetID)
	c.PriceSet = &model.PriceSet{ID: priceSetID}
	if err = c.PriceSet.Read("ID"); err != nil {
		o.Failure("price_set.invalid", util.ErrorInvalidData("price set"))
	}

	subDistrictID, _ = common.Decrypt(c.SubDistrictID)
	c.SubDistrict = &model.SubDistrict{ID: subDistrictID}
	if err = c.SubDistrict.Read("ID"); err != nil {
		o.Failure("sub_district.invalid", util.ErrorInvalidData("sub district"))
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number", util.ErrorCharLength("phone number", 8))
	}

	warehouseID, _ := common.Decrypt(c.WarehouseID)
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if err = c.Warehouse.Read("ID"); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse.Status != 1 {
			o.Failure("warehouse.inactive", util.ErrorActive("default warehouse"))
		}

		filter := map[string]interface{}{"sub_district_id": subDistrictID}
		exclude := map[string]interface{}{}
		if _, countWarehouse, err := repository.CheckWarehouseCoverageData(filter, exclude); err != nil || countWarehouse == 0 {
			o.Failure("warehouse.invalid", util.ErrorMustBeSame("warehouse sub district", "sub district"))
		}
	}

	// Add price set change records into the note
	if c.Branch.PriceSet.ID != c.PriceSet.ID {
		c.NotePriceSetChange = "Priceset Changed | Before: " + c.Branch.PriceSet.Name + " - " + "After: " + c.PriceSet.Name
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"price_set_id.required":     util.ErrorInputRequired("price set"),
		"pic_name.required":         util.ErrorInputRequired("pic name"),
		"phone_number.required":     util.ErrorInputRequired("phone number"),
		"shipping_address.required": util.ErrorInputRequired("shipping address"),
		"sub_district_id.required":  util.ErrorInputRequired("sub district"),
		"warehouse_id.required":     util.ErrorInputRequired("warehouse"),
	}
}
