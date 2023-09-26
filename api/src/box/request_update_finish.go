// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package box

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type updateFinishRequest struct {
	Code string   `json:"-"`
	Rfid []string `json:"rfid" valid:"required"`
	//WarehouseId string `json:"warehouse_id" valid:"required"`
	MacAddress string `json:"-"`

	Session      *auth.SessionData   `json:"-"`
	Box          *model.Box          `json:"-"`
	Warehouse    *model.Warehouse    `json:"-"`
	BoxFridge    *model.BoxFridge    `json:"-"`
	BoxItem      *model.BoxItem      `json:"-"`
	ListBoxItem  []*model.BoxItem    `json:"-"`
	BranchFridge *model.BranchFridge `json:"-"`
}

func (c *updateFinishRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	// WarehouseId, e := common.Decrypt(c.WarehouseId)
	// if e != nil {
	// 	fmt.Println(e)
	// 	o.Failure("tes", e.Error())
	// 	return o
	// }
	// c.Warehouse = &model.Warehouse{ID: WarehouseId}
	// if e = c.Warehouse.Read("ID"); e != nil {
	// 	o.Failure("warehouse.id", e.Error())
	// 	fmt.Println(e)
	// }

	// if e := orSelect.Raw("select * from "+
	// 	"branch_fridge bf "+
	// 	"where bf.status=1 and bf.warehouse_id =? and bf.last_seen_at >= NOW()-INTERVAL ? SECOND", WarehouseId, 10).QueryRow(&c.BranchFridge); e != nil {
	// 	o.Failure("warehouse.status", "fridge offline")
	// }
	if len(c.Rfid) > 0 {
		for _, rfid := range c.Rfid {
			box, e := repository.ValidRfid(rfid)
			if e != nil {
				o.Failure("box", e.Error())
				return o
			}

			c.BoxItem = &model.BoxItem{Box: box, Status: 1}
			if e = c.BoxItem.Read("Box", "Status"); e != nil {
				o.Failure("box_item", e.Error())
			}
			c.ListBoxItem = append(c.ListBoxItem, c.BoxItem)
		}
	}
	// BoxFridge := &model.BoxFridge{BoxItem: c.BoxItem, Status: 1}
	// if e = BoxFridge.Read("Box", "Status"); e != nil {
	// 	o.Failure("tes", e.Error())
	// }

	return o
}

func (c *updateFinishRequest) Messages() map[string]string {
	messages := map[string]string{
		"unit_price.required": util.ErrorInputRequired("Unit Price"),
	}

	return messages
}
