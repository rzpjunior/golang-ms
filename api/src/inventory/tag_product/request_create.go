// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tag_product

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	Name  string   `json:"name" valid:"required"`
	Area  []string `json:"area" valid:"required"`
	Code  string   `json:"-"`
	Image string   `json:"image" valid:"required"`
	Note  string   `json:"note"`
	Value string   `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.Code, err = util.CheckTable("tag_product"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if len(c.Name) < 1 || len(c.Name) > 20 {
		o.Failure("name.invalid", util.ErrorCharLength("name", 20))
	}

	c.Value = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(c.Name, ",", ""), " ", "_"))

	var count int
	orSelect.Raw("select count(tp.id) from tag_product tp where name = ? or value = ?", c.Name, c.Value).QueryRow(&count)
	if count > 0 {
		o.Failure("name.invalid", util.ErrorUnique("name or value"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"name.required":  util.ErrorInputRequired("Name"),
		"area.required":  util.ErrorInputRequired("Area"),
		"image.required": util.ErrorInputRequired("Image"),
	}

	return messages
}
