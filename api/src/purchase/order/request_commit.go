// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type commitRequest struct {
	ID   int64  `json:"-"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (r *commitRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if r.PurchaseOrder.Status != 5 {
		o.Failure("status.inactive", util.ErrorDraft("purchase order"))
		return o
	}

	return o
}

func (r *commitRequest) Messages() map[string]string {
	return map[string]string{}
}
