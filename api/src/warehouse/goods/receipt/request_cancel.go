// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type cancelRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	GoodsReceipt *model.GoodsReceipt `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	r.GoodsReceipt = &model.GoodsReceipt{ID: r.ID}
	r.GoodsReceipt.Read("ID")
	if r.GoodsReceipt.Status != 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("goods receipt", "active"))
	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
