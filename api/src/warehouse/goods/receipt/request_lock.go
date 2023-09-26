// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package receipt

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// lockRequest : struct to lock goods receipt
type lockRequest struct {
	ID        int64 `json:"-" valid:"required"`
	CancelReq int8  `json:"-"`

	GoodsReceipt *model.GoodsReceipt `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate goods receipt request data
func (r *lockRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.GoodsReceipt, err = repository.ValidGoodsReceipt(r.ID); err != nil {
		o.Failure("goods_receipt.invalid", util.ErrorInvalidData("goods receipt"))
		return o
	}

	// To unlock when press cancel in
	if r.GoodsReceipt.Locked == 1 {
		r.CancelReq = 1
	}

	return o
}

// Messages : function to return error validation messages
func (r *lockRequest) Messages() map[string]string {
	messages := map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}

	return messages
}
