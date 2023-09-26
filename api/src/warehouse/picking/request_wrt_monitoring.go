// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// wrtMonitoringRequest : struct to hold wrt monitoring assign request data
type wrtMonitoringRequest struct {
	HelperType string   `json:"helper_type" valid:"required"`
	HelperID   []string `json:"helper_id"`
	WRTID      string   `json:"wrt_id"`

	Staff   []*model.Staff
	WRT     *model.Wrt
	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate wrt monitoring  assign request data
func (r *wrtMonitoringRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if len(r.HelperID) != 0 {
		for _, v := range r.HelperID {
			hID, _ := common.Decrypt(v)

			var staff *model.Staff
			if staff, err = repository.ValidStaff(hID); err != nil {
				o.Failure("helper_id.invalid", util.ErrorInvalidData("staff"))
			}
			r.Staff = append(r.Staff, staff)
		}
	}

	var wrtID int64
	if r.WRTID != "" {
		if wrtID, err = common.Decrypt(r.WRTID); err != nil {
			o.Failure("wrt_id.invalid", util.ErrorInvalidData("WRT"))
		}
		if r.WRT, err = repository.ValidWrt(wrtID); err != nil {
			o.Failure("wrt_id.invalid", util.ErrorInvalidData("WRT"))
		}
	}
	return o
}

// Messages : function to return error validation messages
func (r *wrtMonitoringRequest) Messages() map[string]string {
	return map[string]string{
		"helper_type.required": util.ErrorInputRequired("helper type"),
	}
}
