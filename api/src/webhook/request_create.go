// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package webhook

import (
	"git.edenfarm.id/cuxs/validation"
)

// createRequest : struct to hold sales order request data
type createRequest struct {
	Data map[string]interface{} `json:"data"`
}

// Validate : function to validate sales order request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
