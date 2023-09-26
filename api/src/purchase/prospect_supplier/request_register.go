// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package prospect_supplier

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type registerRequest struct {
	Code				string `json:"-"`
	SupplierTypeID		string `json:"supplier_type_id"`

	ID              	int64  `json:"-" valid:"required"`
	Name            	string `json:"name" valid:"required"`
	PhoneNumber     	string `json:"phone_number" valid:"required"`
	AltPhoneNumber  	string `json:"alt_phone_number"`
	Address         	string `json:"address" valid:"required"`
	PicName         	string `json:"pic_name" valid:"required"`
	PicPhoneNumber     	string `json:"pic_phone_number" valid:"required"`
	SubDistrictID		string `json:"sub_district_id" valid:"required"`
	Note				string `json:"note"`

	SubDistrict 		*model.SubDistrict `json:"-"`
	SupplierType 		*model.SupplierType `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *registerRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	prospectsupplier := &model.ProspectSupplier{Name: c.Name}
	if err = prospectsupplier.Read("Name"); err == nil && prospectsupplier.ID != c.ID {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	if c.SupplierTypeID != "" {
		SupplierTypeID, _ := common.Decrypt(c.SupplierTypeID)
		c.SupplierType = &model.SupplierType{ID: SupplierTypeID}
		c.SupplierType.Read()
	}

	if c.SubDistrictID != "" {
		SubDistrictID, _ := common.Decrypt(c.SubDistrictID)
		c.SubDistrict = &model.SubDistrict{ID: SubDistrictID}
		c.SubDistrict.Read()
	}

	return o
}

// Messages : function to return error validation messages
func (c *registerRequest) Messages() map[string]string {
	return map[string]string{

	}
}
