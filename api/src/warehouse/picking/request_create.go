// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"time"

	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold picking request data
type createRequest struct {
	Code                string    `json:"-"`
	WarehouseID         string    `json:"warehouse_id" valid:"required"`
	RecognitionDate     string    `json:"recognition_date" valid:"required"`
	Note                string    `json:"note"`
	RecognitionDateTime time.Time `json:"-"`

	Warehouse          *model.Warehouse `json:"-"`
	PickingOrderAssign []*itemAssign    `json:"picking_order_assign" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type itemAssign struct {
	ID           string `json:"id"`
	HelperID     string `json:"helper_id" valid:"required"`
	SalesOrderID string `json:"sales_order_id" valid:"required"`

	SalesOrder *model.SalesOrder `json:"-"`
	Helper     *model.Staff
	// data for picking item entry
	SalesOrderItem []*model.SalesOrderItem `json:"-"`
}

// Validate : function to validate picking request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	db := orm.NewOrm()
	db.Using("read_only")
	var e error

	warehouseID, _ := common.Decrypt(c.WarehouseID)
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	c.Warehouse.Read("ID")

	if c.RecognitionDateTime, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInputRequired("recognition date"))
	}

	for n, row := range c.PickingOrderAssign {

		//condition for n SO n Helper
		hID, _ := common.Decrypt(row.HelperID)
		h := &model.Staff{ID: hID}
		if e = h.Read("ID"); e == nil {
			if h.Status != 1 {
				o.Failure("helper.active", util.ErrorActive("helper"))
				return o
			}
			row.Helper = h

		} else {
			o.Failure("helper.invalid", util.ErrorInvalidData("helper"))
			return o
		}

		salesOrderID, _ := common.Decrypt(row.SalesOrderID)
		row.SalesOrder = &model.SalesOrder{ID: salesOrderID}
		if e = row.SalesOrder.Read("ID"); e == nil {
			var countSO int8
			db.Raw("select count(*) from sales_order so where so.id = ? and so.status in (1,9,12)", row.SalesOrder.ID).QueryRow(&countSO)

			if countSO == 0 {
				o.Failure("sales_order"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("sales order"))
				return o
			}

			db.Raw("select * from sales_order_item soi where soi.sales_order_id = ?", row.SalesOrder.ID).QueryRows(&row.SalesOrderItem)
		}

	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":         util.ErrorInputRequired("warehouse"),
		"recognition_date.required":     util.ErrorInputRequired("recognition date"),
		"picking_order_assign.required": util.ErrorSalesOrderCannotBeEmpty(),
	}

	for i := range c.PickingOrderAssign {
		messages["item."+strconv.Itoa(i)+".helper_id.required"] = util.ErrorSelectRequired("helper")
		messages["item."+strconv.Itoa(i)+".sales_order_id.required"] = util.ErrorInputRequired("sales order")

	}

	return messages
}
