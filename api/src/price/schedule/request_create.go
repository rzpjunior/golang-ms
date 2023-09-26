// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schedule

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	InsertProductPrice []*requestItem `json:"prices" valid:"required"`
	PriceSetID         string         `json:"price_set_id" valid:"required"`
	ScheduleDateStr    string         `json:"schedule_date" valid:"required"`
	ScheduleTimeStr    string         `json:"schedule_time" valid:"required"`
	ScheduleDate       time.Time
	ScheduleTime       time.Time

	Session  *auth.SessionData `json:"-"`
	PriceSet *model.PriceSet   `json:"-"`
	Price    *model.Price      `json:"-"`
}

type requestItem struct {
	UnitPrice    float64 `json:"unit_price" valid:"required"`
	ProductID    string  `json:"product_id" valid:"required"`
	PriceSetID   string  `json:"price_set_id"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Salability   int     `json:"salability"`

	PriceSet *model.PriceSet
	Product  *model.Product
	Category *model.Category
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	productCheck := make(map[int64]int8)
	var err error
	var exist int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.ScheduleDate, err = time.Parse("2006-01-02", c.ScheduleDateStr); err != nil {
		o.Failure("schedule_date.invalid", util.ErrorInvalidData("schedule date"))
	}

	priceSetId, _ := common.Decrypt(c.PriceSetID)
	c.PriceSet = &model.PriceSet{ID: priceSetId}
	c.PriceSet.Read("ID")

	orSelect.Raw("SELECT COUNT(id) from price_schedule where price_set_id = ? AND schedule_date = ? AND schedule_time = ? AND status = 1", priceSetId, c.ScheduleDateStr, c.ScheduleTimeStr).QueryRow(&exist)
	if exist > 0 {
		o.Failure("schedule_time.invalid", util.ErrorCreateDoc("scheduler", "schedule"))
	}

	for _, v := range c.InsertProductPrice {
		productId, _ := common.Decrypt(v.ProductID)
		categoryId, _ := common.Decrypt(v.CategoryID)
		v.Product = &model.Product{ID: productId}

		err = v.Product.Read("ID")

		// validation for directory database
		if err != nil {
			o.Failure("id.invalid", util.ErrorMustExistInDirectory("product"))
			return o
		}

		// validation category based on selected price
		if v.CategoryID != "" {
			if v.Product.Category.ID != categoryId {
				o.Failure("id.invalid", util.ErrorMustBeSame("category name", "selected product"))
				return o
			}
		}

		// validation for salability
		if v.Salability != 0 {

			//salability from product itself
			if v.Product.Salability != int8(v.Salability) {
				o.Failure("id.invalid", util.ErrorMustBeSame("salability", "selected product"))
				return o
			}
		}

		// validation for unit price gte 0
		if v.UnitPrice < 0 {
			o.Failure("id.invalid", util.ErrorEqualGreater("unit price", "0"))
			return o
		}

		// validation for duplicate product id
		if _, isExist := productCheck[productId]; isExist {
			o.Failure("id.invalid", util.ErrorDuplicateID("product"))
			return o
		} else {
			productCheck[productId] = 1
		}

	}
	return o
}

func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"schedule_date.required": util.ErrorInputRequired("Schedule Date"),
		"schedule_time.required": util.ErrorInputRequired("Schedule Time"),
		"unit_price.required":    util.ErrorInputRequired("Unit Price"),
	}

	return messages
}
