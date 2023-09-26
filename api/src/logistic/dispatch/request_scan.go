// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dispatch

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strings"
)

// scanRequest : struct to hold dispatch request data
type scanRequest struct {
	SalesOrderCode string `json:"sales_order_code"`
	Increment      string `json:"-"`
	TypeRequest    string `json:"type_request"`

	DeliveryKoliIncrement *model.DeliveryKoliIncrement `json:"-"`
	PickingOrderAssign    *model.PickingOrderAssign    `json:"-"`
	Session               *auth.SessionData            `json:"-"`
}

// Validate : function to validate dispatch request data
func (u *scanRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")

	var poaID int64
	var e error

	if u.TypeRequest == "scan" {
		strArr := strings.Split(u.SalesOrderCode, "-")
		var salesOrder string
		num := strArr[len(strArr)-1]
		u.Increment = num
		salesOrder = strings.TrimSuffix(u.SalesOrderCode, "-"+num)
		u.SalesOrderCode = salesOrder
	}

	if u.SalesOrderCode != "" {
		o1.Raw("select poa.id from sales_order so join picking_order_assign poa on so.id = poa.sales_order_id where so.code = ?", u.SalesOrderCode).QueryRow(&poaID)

		if u.PickingOrderAssign, e = repository.ValidPickingOrderAssign(poaID); e != nil {
			o.Failure("picking_order_id.invalid", util.ErrorInvalidData("picking order"))
			return o
		}

		if u.TypeRequest == "scan" && u.PickingOrderAssign.DispatchStatus == 2 {
			o.Failure("picking_order_id.invalid", util.ErrorActive("dispatch"))
			return o
		}
	}

	if u.TypeRequest == "scan" {
		o1.Raw("SELECT id, sales_order_id, `increment`, is_read, print_label "+
			"FROM eden_v2.delivery_koli_increment "+
			"WHERE `increment` = ? and sales_order_id = ?", u.Increment, u.PickingOrderAssign.SalesOrder.ID).QueryRow(&u.DeliveryKoliIncrement)
	}
	if u.TypeRequest == "scan" && u.DeliveryKoliIncrement.IsRead == 1 {
		o.Failure("delivery_koli_increment.invalid", util.ErrorAlreadyScanned("delivery koli increment"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *scanRequest) Messages() map[string]string {
	messages := map[string]string{}
	return messages
}
