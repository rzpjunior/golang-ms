// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

// updateCommitedRequest : struct to hold request data
type updateCommitedRequest struct {
	Session *auth.SessionData
}

// Validate : function to validate request data
func (c *updateCommitedRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (c *updateCommitedRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
