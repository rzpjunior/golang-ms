// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_deliver

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// printConsolidateRequest : struct to hold print request data
type printConsolidateRequest struct {
	ID int64 `json:"-"`

	ConsolidatedPurchaseDeliver *model.ConsolidatedPurchaseDeliver `json:"-"`
	Session                     *auth.SessionData                  `json:"-"`
}

// Validate : function to validate print request data
func (r *printConsolidateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	r.ConsolidatedPurchaseDeliver = &model.ConsolidatedPurchaseDeliver{ID: r.ID}
	if err := r.ConsolidatedPurchaseDeliver.Read("ID"); err != nil {
		o.Failure("consolidated_purchase_deliver_id.invalid", util.ErrorInvalidData("consolidated purchase deliver id"))
	}

	return o
}

func (r *printConsolidateRequest) Messages() map[string]string {
	return map[string]string{}
}
