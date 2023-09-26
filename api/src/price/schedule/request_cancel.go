// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schedule

import (
	"fmt"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"time"
)

// cancelRequest : struct to hold price schedule request data
type cancelRequest struct {
	ID               int64  `json:"-" valid:"required"`
	CancellationNote string `json:"note" valid:"required"`

	PriceSchedule *model.PriceSchedule
	Session       *auth.SessionData
}

// Validate : function to validate price schedule request data
func (c *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	// current + 1 hour to validate cancellation
	current := time.Now().Add(time.Hour * 1)

	c.PriceSchedule = &model.PriceSchedule{ID: c.ID}
	if err = c.PriceSchedule.Read("ID"); err == nil {
		if c.PriceSchedule.Status != 1 {
			o.Failure("status.inactive", util.ErrorActive("price schedule"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("price schedule"))
	}

	scheduleTime, _ := time.Parse("15:04", c.PriceSchedule.ScheduleTime)
	currentDate := current.Format("2006-01-02")
	currentTime, _ := time.Parse("15:04", fmt.Sprintf("%02d:%02d", current.Hour(), current.Minute()))

	if currentDate == c.PriceSchedule.ScheduleDate {
		if currentTime.After(scheduleTime) {
			o.Failure("id.invalid", "Cannot cancel this scheduler")
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
