// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fridge

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold courier set request data
type createRequestUser struct {
	// User creation
	Username        string `json:"username" valid:"required"`
	Password        string `json:"password" valid:"required"`
	ConfirmPassword string `json:"confirm_password" valid:"required"`
	BranchID        string `json:"branch_id" valid:"required"`
	WarehouseID     string `json:"warehouse_id" valid:"required"`

	PasswordHash string `json:"-"`

	//Branch Creation
	Code      string            `json:"-"`
	Branch    *model.Branch     `json:"-"`
	Warehouse *model.Warehouse  `json:"-"`
	Session   *auth.SessionData `json:"-"`
}

// Validate : function to validate courier request data
func (c *createRequestUser) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	var err error
	lenString := len(c.Username)
	c.Username = c.Username[0:lenString]
	if lenString > 100 {
		c.Username = c.Username[0:100]
	}

	// duplicate user check
	user := &model.UserFridge{Username: c.Username}
	if err = user.Read("Username"); err == nil {
		o.Failure("Username.invalid", util.ErrorDuplicate("Username"))
	}
	// password check and hash
	if errors := util.CheckPassword(c.Password); errors != "" {
		o.Failure("password.invalid", errors)
	}
	if c.ConfirmPassword != c.Password {
		o.Failure("confirm_password.notmatch", "password not match")
	}
	if c.PasswordHash, err = common.PasswordHasher(c.Password); err != nil {
		o.Failure("password.invalid", util.ErrorInvalidData("password"))
	}

	// user fridge validation
	if c.Code, err = util.CheckTable("user_fridge"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	branchID, err := common.Decrypt(c.BranchID)
	if err != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	c.Branch = &model.Branch{ID: branchID}
	if e := c.Branch.Read("ID"); e != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
	}

	warehouseID, err := common.Decrypt(c.WarehouseID)
	if err != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	c.Warehouse = &model.Warehouse{ID: warehouseID}
	if e := c.Warehouse.Read("ID"); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequestUser) Messages() map[string]string {
	messages := map[string]string{
		"username.required":         util.ErrorInputRequired("username"),
		"password.required":         util.ErrorInputRequired("password"),
		"confirm_password.required": util.ErrorInputRequired("confirm password"),
	}

	return messages
}
