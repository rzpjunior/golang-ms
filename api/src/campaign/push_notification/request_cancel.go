// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package push_notification

import (
	"time"

	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// cancelRequest : struct to hold Cancel push notification request data
type cancelRequest struct {
	ID   int64  `json:"-"`
	Note string `json:"note" valid:"required"`

	NotificationCampaign *model.NotificationCampaign `json:"-"`
	Session              *auth.SessionData           `json:"-"`
}

// Validate : function to validate cancel push notification request data
func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	c.NotificationCampaign = &model.NotificationCampaign{ID: c.ID}
	if e = c.NotificationCampaign.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("push notification"))
		return o
	}

	// checking status to glossary
	var statusGlossary *model.Glossary
	if statusGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "status", "value_int", c.NotificationCampaign.Status); e != nil {
		o.Failure("status.invalid", util.ErrorInvalidData("status"))
		return o
	}

	// check status
	if statusGlossary.ValueName != "draft" {
		o.Failure("status."+statusGlossary.ValueName, util.ErrorDocStatus("status", statusGlossary.ValueName))
		return o
	}

	// check greater than time now
	if c.NotificationCampaign.ScheduledAt.Before(time.Now()) {
		o.Failure("note.invalid", util.ErrorGreater("scheduled at", "time now"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required": util.ErrorInputRequired("note"),
	}

	return messages
}
