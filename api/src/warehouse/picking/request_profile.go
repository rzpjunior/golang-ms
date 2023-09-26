// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

// profileRequestPicking : struct to hold picking assign request data
type profileRequestPicking struct {
	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *profileRequestPicking) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (r *profileRequestPicking) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
