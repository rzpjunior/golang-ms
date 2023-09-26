// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package _return

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type cancelRequest struct {
	ID               int64  `json:"-"`
	CancellationNote string `json:"note" valid:"required"`

	DeliveryReturn *model.DeliveryReturn `json:"-"`
	Session        *auth.SessionData     `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	r.DeliveryReturn = &model.DeliveryReturn{ID: r.ID}
	if err = r.DeliveryReturn.Read("ID"); err == nil {
		if r.DeliveryReturn.Status != 1 {
			o.Failure("status.inactive", util.ErrorActive("delivery return"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("delivery return"))
	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("note"),
	}
}
