// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package plan

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type confirmRequest struct {
	ID int64 `json:"-" valid:"required"`

	PurchasePlan *model.PurchasePlan `json:"-"`
	Session      *auth.SessionData   `json:"-"`
}

func (r *confirmRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.PurchasePlan, err = repository.ValidPurchasePlan(r.ID)
	if err != nil {
		o.Failure("purchase_plan_id.invalid", util.ErrorInvalidData("purchase plan"))
	}

	if r.PurchasePlan.Status != 1 {
		o.Failure("purchase_plan.invalid", util.ErrorActive("purchase plan"))
	}

	return o
}

func (r *confirmRequest) Messages() map[string]string {
	return map[string]string{}
}
