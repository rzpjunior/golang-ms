// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package price_set

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID          int64  `json:"-" valid:"required"`
	Name        string `json:"name" valid:"required"`
	Note    	string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	priceset := &model.PriceSet{Name: c.Name}
	if err = priceset.Read("Name"); err == nil && priceset.ID != c.ID {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required": util.ErrorInputRequired("name"),
	}
}
