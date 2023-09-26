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

// createShippingAddressRequest : struct to hold price set request data
type createShippingAddressRequest struct {
	Code string `json:"-"`
	//AddressName             string `json:"address_name" valid:"required"`
	RecipientName           string `json:"recipient_name" valid:"required"`
	RecipientPhoneNumber    string `json:"recipient_phone_number" valid:"required"`
	RecipientAltPhoneNumber string `json:"recipient_alt_phone_number"`
	ShippingAddress         string `json:"shipping_address" valid:"required"`
	ShippingNote            string `json:"shipping_note"`
	MerchantId              string `json:"merchant_id" valid:"required"`
	ArchetypeId             string `json:"archetype_id" valid:"required"`
	SalespersonId           string `json:"salesperson_id" valid:"required"`
	AreaId                  string `json:"area_id" valid:"required"`
	SubDistrictId           string `json:"sub_district_id" valid:"required"`

	Merchant  *model.Merchant
	Archetype *model.Archetype
	//PriceSet          *model.PriceSet
	Salesperson       *model.Staff
	Area              *model.Area
	SubDistrict       *model.SubDistrict
	Warehouse         *model.Warehouse
	WarehouseCoverage *model.WarehouseCoverage

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *createShippingAddressRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Code, err = util.CheckTable("branch"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if merchantId, err := common.Decrypt(c.MerchantId); err != nil {
		o.Failure("merchant_id.invalid", util.ErrorInvalidData("merchant"))
	} else {
		if c.Merchant, err = repository.ValidMerchant(merchantId); err != nil {
			o.Failure("merchant_id.invalid", util.ErrorInvalidData("merchant"))
		} else {
			if c.Merchant.Status != int8(1) {
				o.Failure("merchant_id.active", util.ErrorActive("merchant"))
			}
		}
	}

	if archetypeId, err := common.Decrypt(c.ArchetypeId); err != nil {
		o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
	} else {
		if c.Archetype, err = repository.ValidArchetype(archetypeId); err != nil {
			o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
		} else {
			if c.Archetype.Status != int8(1) {
				o.Failure("archetype_id.active", util.ErrorActive("archetype"))
			}
		}
	}

	//if priceSetId, err := common.Decrypt(c.PriceSetId); err != nil {
	//	o.Failure("price_set_id.invalid", util.ErrorInvalidData("price_set"))
	//} else {
	//	if c.PriceSet, err = repository.ValidPriceSet(priceSetId); err != nil {
	//		o.Failure("price_set_id.invalid", util.ErrorInvalidData("price_set"))
	//	} else {
	//		if c.PriceSet.Status != int8(1) {
	//			o.Failure("price_set_id.active", util.ErrorActive("price_set"))
	//		}
	//	}
	//}

	if salespersonId, err := common.Decrypt(c.SalespersonId); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales person"))
	} else {
		if c.Salesperson, err = repository.ValidStaff(salespersonId); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("sales person"))
		} else {
			if c.Salesperson.Status != int8(1) {
				o.Failure("id.invalid", util.ErrorActive("sales person"))
			}
		}
	}

	if areaId, err := common.Decrypt(c.AreaId); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("shipping area"))
	} else {
		if c.Area, err = repository.ValidArea(areaId); err != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("shipping area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area_id.active", util.ErrorActive("shipping area"))
			}
		}
	}

	if subDistrictId, err := common.Decrypt(c.SubDistrictId); err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	} else {
		if c.SubDistrict, err = repository.ValidSubDistrict(subDistrictId); err != nil {
			o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
		}
		if c.SubDistrict.Area.ID != c.Area.ID {
			o.Failure("sub_district_id.invalid", util.ErrorMustBeSame("sub district", "area"))
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
		o.Failure("warehouse_coverage.invalid", util.ErrorMustBeSame("warehouse coverage", "warehouse"))
	}

	if len(c.ShippingNote) > 100 {
		o.Failure("shipping_note.invalid", util.ErrorCharLength("shipping note", 100))
	}

	if len(c.ShippingAddress) > 350 {
		o.Failure("shipping_address.invalid", util.ErrorCharLength("shipping address", 350))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createShippingAddressRequest) Messages() map[string]string {
	return map[string]string{
		"merchant_id.required":            util.ErrorInputRequired("merchant id"),
		"recipient_name.required":         util.ErrorInputRequired("recipient name"),
		"recipient_phone_number.required": util.ErrorInputRequired("recipient phone number"),
		"shipping_address.required":       util.ErrorInputRequired("shipping address"),
		"area_id.required":                util.ErrorInputRequired("shipping area"),
		"sub_district_id.required":        util.ErrorInputRequired("sub district"),
		"default_warehouse_id.required":   util.ErrorInputRequired("default warehouse"),
	}
}
