// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package day_off

import (
	"time"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold wrt config request data
type createRequest struct {
	OffDateStr string    `json:"off_date" valid:"required"`
	Note       string    `json:"note"`
	OffDate    time.Time `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate wrt config request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	layout := "2006-01-02"

	if c.OffDate, err = time.Parse(layout, c.OffDateStr); err != nil {
		o.Failure("off_date.invalid", util.ErrorInvalidData("date"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"off_date.required": util.ErrorInputRequired("date"),
	}
}
