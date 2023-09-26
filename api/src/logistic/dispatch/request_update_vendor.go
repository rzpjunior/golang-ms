// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dispatch

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// updateVendorRequest : struct to hold dispatch request data
type updateVendorRequest struct {
	ID              int64  `json:"-" valid:"required"`
	CourierVendorID string `json:"courier_vendor_id" valid:"required"`
	CourierID       string `json:"courier_id" valid:"required"`

	CourierVendor      *model.CourierVendor      `json:"-"`
	Courier            *model.Courier            `json:"-"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	Session            *auth.SessionData         `json:"-"`
}

// Validate : function to validate dispatch request data
func (u *updateVendorRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	u.PickingOrderAssign = &model.PickingOrderAssign{ID: u.ID}
	u.PickingOrderAssign.Read("ID")

	if u.PickingOrderAssign.DispatchStatus != 1 {
		o.Failure("courier_vendor_id.inactive", util.ErrorActiveInd("dispatch"))
		return o
	}

	if vID, e := common.Decrypt(u.CourierVendorID); e != nil {
		o.Failure("courier_vendor_id.invalid", util.ErrorInvalidData("courier vendor"))
	} else {
		if u.CourierVendor, e = repository.ValidCourierVendor(vID); e != nil {
			o.Failure("courier_vendor_id.invalid", util.ErrorInvalidData("courier vendor"))
			return o
		}

		if u.CourierVendor.Status != 1 {
			o.Failure("courier_vendor.inactive", util.ErrorActiveInd("courier vendor"))
			return o
		}
	}
	if cID, e := common.Decrypt(u.CourierID); e != nil {
		o.Failure("courier_id.invalid", util.ErrorInvalidData("courier"))
	} else {
		if u.Courier, e = repository.ValidCourier(cID); e != nil {
			o.Failure("courier_id.invalid", util.ErrorInvalidData("courier"))
			return o
		}

		if u.Courier.Status != 1 {
			o.Failure("courier.inactive", util.ErrorActiveInd("courier"))
			return o
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateVendorRequest) Messages() map[string]string {
	messages := map[string]string{
		"courier_id.required":        util.ErrorInputRequired("courier"),
		"courier_vendor_id.required": util.ErrorInputRequired("courier vendor"),
	}
	return messages
}
