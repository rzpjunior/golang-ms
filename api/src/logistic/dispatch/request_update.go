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

// updateRequest : struct to hold dispatch request data
type updateRequest struct {
	ID        int64  `json:"-" valid:"required"`
	CourierID string `json:"courier_id"`

	Courier            *model.Courier            `json:"-"`
	PickingOrderAssign *model.PickingOrderAssign `json:"-"`
	Session            *auth.SessionData         `json:"-"`
}

// Validate : function to validate dispatch request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	u.PickingOrderAssign = &model.PickingOrderAssign{ID: u.ID}
	u.PickingOrderAssign.Read("ID")

	if u.CourierID == "" {
		u.Courier = nil
	} else {
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

	}
	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"courier_id.required": util.ErrorInputRequired("courier"),
	}
	return messages
}
