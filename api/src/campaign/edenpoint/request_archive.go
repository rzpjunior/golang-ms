// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenpoint

import (
	"time"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// archiveRequest : struct to hold request data
type archiveRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	EdenPointCampaign *model.EdenPointCampaign `json:"-"`

	Session *auth.SessionData
}

// Validate : function to validate request data
func (r *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if r.EdenPointCampaign, err = repository.ValidEdenPointCampaign(r.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("id"))
	}

	currentTime := time.Now()
	if currentTime.After(r.EdenPointCampaign.EndDate) {
		o.Failure("id.invalid", util.ErrorActive("id"))
	}

	r.EdenPointCampaign.Status = 3

	return o
}

// Messages : function to return error validation messages
func (r *archiveRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("note"),
	}
}
