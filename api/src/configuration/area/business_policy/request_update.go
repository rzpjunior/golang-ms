// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package business_policy

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateRequest struct {
	ID             int64   `json:"-" valid:"required"`
	AreaID         string  `json:"area_id" valid:"required"`
	BusinessTypeID string  `json:"business_type_id" valid:"required"`
	MinOrder       float64 `json:"min_order" valid:"required"`
	DeliveryFee    float64 `json:"delivery_fee" valid:"required"`

	Area         *model.Area         `json:"-"`
	BusinessType *model.BusinessType `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		e                      error
		areaID, businessTypeID int64
	)

	if areaID, e = common.Decrypt(c.AreaID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		return o
	}

	if c.Area, e = repository.ValidArea(areaID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		return o
	}

	if c.Area.Status != int8(1) {
		o.Failure("area_id.invalid", util.ErrorActive("area"))
		return o
	}

	if businessTypeID, e = common.Decrypt(c.BusinessTypeID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("business type"))
		return o
	}

	if c.BusinessType, e = repository.ValidBusinessType(businessTypeID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("business type"))
		return o
	}

	if c.BusinessType.Status != int8(1) {
		o.Failure("area_id.invalid", util.ErrorActive("business type"))
		return o
	}

	if c.MinOrder == 0 {
		o.Failure("min_order.invalid", util.ErrorInputRequired("min order free delivery"))
	}

	if c.DeliveryFee == 0 {
		o.Failure("delivery_fee.invalid", util.ErrorInputRequired("delivery fee"))
	}

	return o
}

// Messages : function to return error messages after validation
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"area_id.required":          util.ErrorInputRequired("area"),
		"business_type_id.required": util.ErrorInputRequired("business type"),
		"min_order.required":        util.ErrorInputRequired("min order free delivery"),
		"delivery_fee.required":     util.ErrorInputRequired("delivery fee"),
	}
}
