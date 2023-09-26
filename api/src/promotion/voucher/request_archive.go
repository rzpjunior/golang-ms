// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package voucher

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"time"
)

type archiveRequest struct {
	ID         int64             `json:"-" valid:"required"`
	VoidReason int8              `json:"-"`
	Session    *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	if voucher, err := repository.ValidVoucher(c.ID); err == nil {
		if voucher.Status != 1 {
			o.Failure("id.invalid", util.ErrorActive("status"))
		} else {
			if voucher.RemOverallQuota == 0 {
				c.VoidReason = 2
			} else if voucher.EndTimestamp.Before(time.Now()) {
				c.VoidReason = 1
			} else {
				c.VoidReason = 3
			}

		}
	} else {
		o.Failure("voucher.invalid", util.ErrorInvalidData("voucher"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
