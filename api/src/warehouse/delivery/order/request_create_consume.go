// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

// createConsumeRequest : struct to hold consumer request data
type createConsumeRequest struct {
	Note string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *createConsumeRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (c *createConsumeRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
