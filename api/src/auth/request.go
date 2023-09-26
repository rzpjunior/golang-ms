// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/util"
)

// SignInRequest data struct that stored request data when requesting an create auth process.
// All data must be provided and must be match with specification validation below.
// handler function should be bind this with context to matches incoming request
// data keys to the defined json tag.
type SignInRequest struct {
	Email                  string      `json:"email"`
	Password               string      `valid:"required" json:"password"`
	FirebaseTokenDashboard string      `json:"firebase_token_dashboard"`
	IPAddress              string      `json:"-"`
	User                   *model.User `json:"-"`
}

// Validate implement validation.Requests interfaces.
func (r *SignInRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	//to get expired period to limit how many seconds value from database

	configApp, e := repository.GetConfigApp("attribute", "dino_failed_login_retry_time")

	period, _ := strconv.Atoi(configApp.Value)
	expPeriod := time.Duration(period) * time.Second

	counterRedis := 0

	// check email di user
	user := new(model.User)
	user.Email = r.Email

	redisKey := r.IPAddress + "_" + r.Email
	if e = user.Read("Email"); e != nil {
		if !isLimit(redisKey, expPeriod) {
			o.Failure("id", "Too many failed login attempts for this email - Please try again in "+configApp.Value+" seconds")
			return o
		}
		counterRedis = counterRedis + 1

		o.Failure("email", util.ErrInvalidCredential)
	}

	if user.Email != "superadmin" {
		staff := &model.Staff{User: user}
		if e = staff.Read("User"); e != nil {
			o.Failure("email", util.ErrInvalidCredential)
		}
		if e = staff.Role.Read("ID"); e != nil {
			o.Failure("email", util.ErrInvalidCredential)
		}
		if staff.Role.Code == "ROL0022" {
			o.Failure("email", util.ErrInvalidCredential)
		}
	}
	// cek user apakah active?
	if user.Status != 1 {
		o.Failure("email", "Your account is not activated yet")
	}
	// cek password sesuai dengan inputan
	if err := common.PasswordHash(user.Password, r.Password); err != nil {
		if counterRedis == 1 {
			o.Failure("email", util.ErrInvalidCredential)
			return o
		}
		if !isLimit(redisKey, expPeriod) {
			o.Failure("id", "Too many failed login attempts for this email - Please try again in "+configApp.Value+" seconds")
			return o
		}
		o.Failure("email", util.ErrInvalidCredential)
	}

	if !isMaxLimit(redisKey, expPeriod) {
		o.Failure("id", "Too many failed login attempts for this email - Please try again in "+configApp.Value+" seconds")
		return o
	}

	if o.Valid {
		r.User = user
	}

	return o
}

// Messages implement validation.Requests interfaces
// return custom messages when validation fails.
func (r *SignInRequest) Messages() map[string]string {
	return map[string]string{
		"password.invalid": util.ErrorInputRequired("password"),
	}
}
