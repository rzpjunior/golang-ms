// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package box

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

// printRequest : struct to hold price set request data
type printRequest struct {
	RfidCode string            `json:"rfid_code"`
	Session  *auth.SessionData `json:"-"`
}

// Validate : function to validate uom request data
func (c *printRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (c *printRequest) Messages() map[string]string {
	return map[string]string{}
}
