// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package transfer_sku

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// cancelRequest : struct to hold Cancel Transfer SKU request data
type cancelRequest struct {
	ID          int64              `json:"-" valid:"required"`
	TransferSku *model.TransferSku `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate Cancel Transfer SKU request data
func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.TransferSku, err = repository.GetTransferSku("ID", c.ID); err != nil {
		o.Failure("id_invalid", util.ErrorInvalidData("transfer sku"))
		return o
	}

	if err = c.TransferSku.Warehouse.Read("ID"); err != nil {
		o.Failure("id_invalid", util.ErrorInvalidData("warehouse"))
		return o
	}

	if c.TransferSku.Status != 1 {
		o.Failure("id_invalid", util.ErrorActive("transfer sku"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"id.required": util.ErrorInputRequired("id"),
	}
}
