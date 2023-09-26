// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// cancelRequest : struct to hold price set request data
type cancelRequest struct {
	ID int64 `json:"-"`

	PackingOrder *model.PackingOrder `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.PackingOrder, err = repository.ValidPackingOrder(c.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("packing order"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
