// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wrt

import (
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold wrt config request data
type createRequest struct {
	Code   string `json:"-"`
	Name   string `json:"name" valid:"required"`
	Note   string `json:"note"`
	AreaId string `json:"area_id" valid:"required"`
	Type   int8   `json:"type"`

	Area *model.Area

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate wrt config request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("wrt"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	areaId, e := common.Decrypt(c.AreaId)
	if e != nil {
		o.Failure("area.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, e = repository.ValidArea(areaId); e != nil {
			o.Failure("area.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area.invalid", util.ErrorActive("area"))
			}
		}
	}

	filter := map[string]interface{}{"name": c.Name, "area_id": strconv.Itoa(int(areaId))}
	exclude := map[string]interface{}{"status": 3}
	if _, countName, err := repository.CheckWrtData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	// Set default type to 1 (delivery)
	if c.Type == 0 {
		c.Type = 1
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":    util.ErrorInputRequired("name"),
		"area_id.required": util.ErrorSelectRequired("area"),
	}
}
