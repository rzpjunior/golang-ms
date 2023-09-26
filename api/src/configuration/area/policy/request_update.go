// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID                  int64   `json:"-" valid:"required"`
	MinOrder            float64 `json:"min_order" valid:"required"`
	DeliveryFee         float64 `json:"delivery_fee" valid:"required"`
	OrderTimeLimit      string  `json:"order_time_limit" valid:"required"`
	AreaID              string  `json:"area_id" valid:"required"`
	DefaultPriceSetID   string  `json:"default_price_set" valid:"required"`
	MaxDayDeliveryDate  int     `json:"-"`
	WeeklyDayOff        int     `json:"-"`
	DraftOrderTimeLimit string  `json:"-"`

	OrderTimeLimitFr time.Time `json:"-"`

	Area            *model.Area     `json:"-"`
	DefaultPriceSet *model.PriceSet `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	if AreaID, e := common.Decrypt(c.AreaID); e != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, e = repository.ValidArea(AreaID); e != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area_id.invalid", util.ErrorActive("area"))
			}
		}
	}

	if c.MinOrder == 0 {
		o.Failure("min_order.invalid", util.ErrorInputRequired("min order free delivery"))
	}

	if c.DeliveryFee == 0 {
		o.Failure("delivery_fee.invalid", util.ErrorInputRequired("delivery fee"))
	}

	// only for checking format time from apps
	if c.OrderTimeLimitFr, e = time.Parse("15:04", c.OrderTimeLimit); e != nil {
		o.Failure("order_time_limit.invalid", util.ErrorInvalidData("order time limit"))
	}

	if c.DefaultPriceSetID != "" {
		DefaultPriceSetID, _ := common.Decrypt(c.DefaultPriceSetID)
		c.DefaultPriceSet = &model.PriceSet{ID: DefaultPriceSetID}
		c.DefaultPriceSet.Read()
	}

	areaPolicy := &model.AreaPolicy{ID: c.ID}
	areaPolicy.Read("ID")
	c.MaxDayDeliveryDate = areaPolicy.MaxDayDeliveryDate
	c.WeeklyDayOff = areaPolicy.WeeklyDayOff
	c.DraftOrderTimeLimit = areaPolicy.DraftOrderTimeLimit

	return o
}

// Messages : function to return error messages after validation
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"min_order.required":         util.ErrorInputRequired("min order free delivery"),
		"delivery_fee.required":      util.ErrorInputRequired("delivery fee"),
		"order_time_limit.required":  util.ErrorInputRequired("order time limit"),
		"default_price_set.required": util.ErrorInputRequired("price set"),
		"area_id.required":           util.ErrorInputRequired("area"),
	}
}
