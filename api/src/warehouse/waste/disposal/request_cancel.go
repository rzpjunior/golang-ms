// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package disposal

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type cancelRequest struct {
	ID               int64  `json:"-"`
	CancellationNote string `json:"cancellation_note" valid:"required"`

	WasteDisposal *model.WasteDisposal `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.WasteDisposal = &model.WasteDisposal{ID: r.ID}
	if err = r.WasteDisposal.Read("ID"); err == nil {
		if r.WasteDisposal.Status != 1 {
			o.Failure("status.inactive", util.ErrorActive("waste disposal"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("waste disposal"))
	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"cancellation_note.required": util.ErrorInputRequired("cancellation note"),
	}
}
