// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package branch

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// archiveRequest : struct to hold outlet request data
type archiveRequest struct {
	ID              int64 `json:"-" valid:"required"`
	ArchiveMerchant int8

	Branch   *model.Branch
	Merchant *model.Merchant

	Session *auth.SessionData
}

// Validate : function to validate customer tag request data
func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	// variable to set whether merchant should also be archived or not. 1 -> archive merchant & user merchant, 2 -> not archive merchant
	c.ArchiveMerchant = 2

	if c.Branch, err = repository.ValidBranch(c.ID); err == nil {
		if c.Branch.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		} else {
			if countActiveBranch, err := repository.CountActiveBranchByMerchantId(c.Branch.Merchant.ID); err == nil {
				if countActiveBranch == 1 {
					c.ArchiveMerchant = 1
					c.Merchant = &model.Merchant{
						ID: c.Branch.Merchant.ID,
					}
					if err = c.Merchant.Read("ID"); err != nil {
						o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
					}
				}
			}
		}
	} else {
		o.Failure("branch.invalid", util.ErrorInvalidData("branch"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
