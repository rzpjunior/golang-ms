// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// assignLeadPickerRequest : struct to hold picking assign request data
type assignLeadPickerRequest struct {
	ID              int64     `json:"-"`
	HelperID        string    `json:"staff_id"`
	AssignTimeStamp time.Time `json:"-"`

	PickingList *model.PickingList `json:"-"`
	Helper      *model.Staff       `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *assignLeadPickerRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var hID int64

	if r.PickingList, err = repository.ValidPickingList(r.ID); err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}
	if r.PickingList.Status != 1 {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}

	if r.HelperID == "" {
		r.Helper = nil
		r.AssignTimeStamp = time.Time{}
	} else {
		if hID, err = common.Decrypt(r.HelperID); err != nil {
			o.Failure("helper_id.invalid", util.ErrorInvalidData("helper"))
			return o
		}
		if r.Helper, err = repository.ValidStaff(hID); err != nil {
			o.Failure("helper_id.invalid", util.ErrorInvalidData("helper"))
			return o
		}

		if r.Helper.Status != 1 {
			o.Failure("helper_id.inactive", util.ErrorActiveInd("helper"))
			return o
		}
		r.AssignTimeStamp = time.Now()
	}

	return o
}

// Messages : function to return error validation messages
func (r *assignLeadPickerRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
