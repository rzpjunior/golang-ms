// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
)

// SignInRequest data struct that stored request data when requesting an create auth process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type SignInPackingRequest struct {
	Email         string `json:"email"`
	Password      string `valid:"required" json:"password"`
	FirebaseToken string `json:"firebase_token"`

	User *model.User `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *SignInPackingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	// check email di user
	user := new(model.User)
	user.Email = r.Email

	if e = user.Read("Email"); e != nil {
		o.Failure("email", util.ErrInvalidCredential)
	} else {
		// cek user apakah active?
		if user.Status != 1 {
			o.Failure("email", "Your account is not activated yet")
		}
		// cek password sesuai dengan inputan
		if err := common.PasswordHash(user.Password, r.Password); err != nil {
			o.Failure("email", util.ErrInvalidCredential)
		}
	}

	if o.Valid {
		r.User = user
		r.User.PickingNotifToken = r.FirebaseToken
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *SignInPackingRequest) Messages() map[string]string {
	return map[string]string{
		"password.invalid": util.ErrorInputRequired("password"),
	}
}
