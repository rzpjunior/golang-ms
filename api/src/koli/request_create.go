// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package koli

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code  string `json:"-"`
	Value string `json:"value" valid:"value" valid:"required"`
	Name  string `json:"name" valid:"name" valid:"required"`
	Note  string `json:"note" valid:"note"`

	DeliveryOrder      *model.DeliveryOrder      `json:"-"`
	CourierTransaction *model.CourierTransaction `json:"-"`
	Session            *auth.SessionData         `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	var err error
	o := &validation.Output{Valid: true}

	koli := &model.Koli{Name: c.Name}
	if err = koli.Read("Name"); err == nil {
		o.Failure("name.invalid", util.ErrorDuplicate("name"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"value.required": util.ErrorInputRequired("value"),
		"name.required":  util.ErrorInputRequired("name"),
	}

	return messages
}
