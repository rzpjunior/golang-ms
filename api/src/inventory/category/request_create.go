// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package category

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strings"
)

// createRequest : struct to hold category request data
type createRequest struct {
	Code                 string `json:"code" valid:"required"`
	Name                 string `json:"name" valid:"required"`
	Classification       int8   `json:"classification" valid:"required"`
	GrandParentID        string `json:"grand_parent_id"`
	ParentID             string `json:"parent_id"`
	GrandParentIDDecrypt int64  `json:"-"`
	ParentIDDecrypt      int64  `json:"-"`
	Note                 string `json:"note"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate category request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error

	filterCode := map[string]interface{}{"code": c.Code}
	excludeStatus := map[string]interface{}{"status": int8(3)}
	if _, countCode, err := repository.CheckCategoryData(filterCode, excludeStatus); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	} else if countCode > 0 {
		o.Failure("code", util.ErrorDuplicate("code"))
	}

	filter := map[string]interface{}{"name": c.Name}
	exclude := map[string]interface{}{"status": int8(3)}
	if _, countName, err := repository.CheckCategoryData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	// as grand parent
	switch c.Classification {
	// as parent
	case 2:
		if c.GrandParentID == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("grand parent"))
		}
		if c.GrandParentIDDecrypt, err = common.Decrypt(c.GrandParentID); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("grand parent"))
		}
	//as child
	case 3:
		if c.GrandParentID == "" || c.ParentID == "" {
			o.Failure("id.invalid", util.ErrorInputRequired("grand parent and parent"))
		}
		if c.GrandParentIDDecrypt, err = common.Decrypt(c.GrandParentID); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("grand parent"))
		}
		if c.ParentIDDecrypt, err = common.Decrypt(c.ParentID); err != nil {
			o.Failure("id.invalid", util.ErrorInvalidData("parent"))
		}

	}
	if len(strings.TrimSpace(c.Code)) != 8 {
		o.Failure("code.invalid", util.ErrorCharLength("code", 8))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":      util.ErrorInputRequired("name"),
		"image_url.required": util.ErrorInputRequired("image_url"),
	}
}
