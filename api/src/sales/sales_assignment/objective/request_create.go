// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package objective

import (
	"net/url"

	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold task assignment objective request data
type createRequest struct {
	Code       string  `json:"-"`
	Name       string  `json:"name" valid:"required"`
	Objective  string  `json:"objective" valid:"required"`
	SurveyLink *string `json:"surveylink"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate create request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("sales_assignment_objective"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if c.SurveyLink != nil {
		if _, err := url.ParseRequestURI(*c.SurveyLink); err != nil {
			o.Failure("surveylink.invalid", util.ErrorInvalidData("surveylink"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":      util.ErrorInputRequired("name"),
		"objective.required": util.ErrorInputRequired("objective"),
	}
}
