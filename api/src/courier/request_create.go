// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package courier

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	DeliveryOrderCode string `json:"delivery_order_code" valid:"required"`
	Latitude          string `json:"latitude" valid:"required"`
	Longitude         string `json:"longitude" valid:"required"`
	Accuracy          string `json:"accuracy" valid:"required"`
	CourierName       string `json:"courier_name" valid:"required"`
	CourierPhoneNo    string `json:"courier_phone_no"`
	Note              string `json:"note"`

	DeliveryOrder      *model.DeliveryOrder      `json:"-"`
	CourierTransaction *model.CourierTransaction `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	var err error
	o := &validation.Output{Valid: true}

	c.DeliveryOrder = &model.DeliveryOrder{Code: c.DeliveryOrderCode}
	if err = c.DeliveryOrder.Read("Code"); err == nil {
		if c.DeliveryOrder.Status != 1 {
			o.Failure("delivery_order_code.inactive", util.ErrorNotFound("delivery order code"))
			return o
		} else {
			orSelect := orm.NewOrm()
			orSelect.Using("read_only")
			orSelect.Raw("select * from courier_transaction ct where ct.delivery_order_id = ?", c.DeliveryOrder.ID).QueryRow(&c.CourierTransaction)
		}
	} else {
		o.Failure("delivery_order_code.invalid", util.ErrorNotFound("delivery order code"))
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"delivery_order_code.required": util.ErrorInputRequired("delivery order code"),
		"latitude.required":            util.ErrorInputRequired("latitude"),
		"longitude.required":           util.ErrorInputRequired("longitude"),
		"accuracy.required":            util.ErrorInputRequired("accuracy"),
		"courier_name.required":        util.ErrorInputRequired("courier name"),
	}

	return messages
}
