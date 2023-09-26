// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consolidated_shipment

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// printRequest : struct to hold print request data
type printRequest struct {
	ID int64 `json:"-"`

	ConsolidatedShipment *model.ConsolidatedShipment `json:"-"`
	Session              *auth.SessionData           `json:"-"`
}

// Validate : function to validate print request data
func (r *printRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	r.ConsolidatedShipment = &model.ConsolidatedShipment{ID: r.ID}
	if err := r.ConsolidatedShipment.Read("ID"); err != nil {
		o.Failure("consolidated_shipment_id.invalid", util.ErrorInvalidData("consolidated shipment id"))
	}

	return o
}

func (r *printRequest) Messages() map[string]string {
	return map[string]string{}
}
