// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms_viro

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateRequest struct {
	Results     []results          `json:"results"`
	OTPOutGoing *model.OtpOutgoing `json:"-"`
}
type results struct {
	Status    status `json:"status"`
	MessageId string `json:"messageId"`
}
type status struct {
	GroupID int `json:"groupId"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	c.OTPOutGoing = &model.OtpOutgoing{VendorMessageID: c.Results[0].MessageId}
	if e = c.OTPOutGoing.Read("VendorMessageID"); e != nil {
		o.Failure("messageId", e.Error())
	}
	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{}
}
