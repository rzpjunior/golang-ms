// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sku_discount

import (
	"time"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// archiveRequest : struct to hold request data
type archiveRequest struct {
	ID int64 `json:"-"`

	SkuDiscount *model.SkuDiscount `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (r *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err         error
		currentTime time.Time
	)

	if r.SkuDiscount, err = repository.ValidSkuDiscount(r.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("discount"))
		return o
	}

	if r.SkuDiscount.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("discount"))
		return o
	}

	if currentTime = time.Now(); currentTime.After(r.SkuDiscount.EndTimestamp) {
		o.Failure("id.invalid", util.ErrorOutOfPeriod("discount"))
		return o
	}

	r.SkuDiscount.Status = 2

	return o
}

// Messages : function to return error validation messages
func (r *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
