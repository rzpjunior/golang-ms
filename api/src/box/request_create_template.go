// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package box

import (
	"math"
	"strconv"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type createRequestTemplate struct {
	WarehouseId     string            `json:"warehouse_id,omitempty"`
	ProductBoxList  []*taskProductBox `json:"productBoxes" valid:"required"`
	NewProductBox   []*taskProductBox `json:"-"`
	ExistProductBox []*taskProductBox `json:"-"`
	listRFID        []string

	Session      *auth.SessionData   `json:"-"`
	Warehouse    *model.Warehouse    `json:"-"`
	BranchFridge *model.BranchFridge `json:"-"`
}

type taskProductBox struct {
	Rfid        string         `json:"Rfid_Code" `
	ProductCode string         `json:"Product_Code" `
	TotalWeight float64        `json:"Total_Weight"`
	UnitPrice   float64        `json:"Unit_Price,omitempty" `
	TotalPrice  float64        `json:"Total_Price,omitempty" `
	Size        float64        `json:"Size_Box,omitempty" `
	Note        string         `json:"Note"`
	Product     *model.Product `json:"-"`
	Box         *model.Box     `json:"-"`
}

func (c *createRequestTemplate) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	// WarehouseId, e := common.Decrypt(c.WarehouseId)
	// if e != nil {
	// 	fmt.Println(e)
	// 	o.Failure("tes", e.Error())
	// }

	// c.Warehouse = &model.Warehouse{ID: WarehouseId}
	// if e := c.Warehouse.Read("ID"); e != nil {
	// 	o.Failure("warehouse.id", e.Error())
	// 	fmt.Println(e)
	// }
	for i, v := range c.ProductBoxList {
		// if v.UnitPrice <= 0 {
		// 	o.Failure("unit_price_."+strconv.Itoa(i)+".invalid", "unit price must be higher than 0")
		// }
		if v.Rfid == "" {
			o.Failure("rfid_code_."+strconv.Itoa(i)+".invalid", "rfid must be not empty")
			continue
		}
		if v.ProductCode == "" {
			o.Failure("product_code_."+strconv.Itoa(i)+".invalid", "product code must be not empty")
			continue
		}
		if v.TotalWeight <= 0 {
			o.Failure("total_weight_."+strconv.Itoa(i)+".invalid", "total weight must be more than 0")
		}
		// if v.Size < 1 {
		// 	o.Failure("size_."+strconv.Itoa(i)+".invalid", "")
		// }
		// if v.Size > 4 {
		// 	o.Failure("size_."+strconv.Itoa(i)+".invalid", "")
		// }
		//validasi panjang char rfid
		// if len(v.Rfid) < 24 {
		// 	continue
		// } else if len(v.Rfid) > 24 {
		// 	//ambil 12 char pertama tes
		// 	v.Rfid = v.Rfid[0:24]
		// }

		//cek if its already processed,then continue
		isEPCExist := util.StringInSlice(v.Rfid, c.listRFID)
		if !isEPCExist {
			c.listRFID = append(c.listRFID, v.Rfid)
		} else {
			continue
		}

		//check is it exist in table box based on rfid
		box, _ := repository.ValidRfid(v.Rfid)

		v.Product = &model.Product{Code: v.ProductCode}
		if e := v.Product.Read("Code"); e != nil {
			o.Failure("product_code_"+strconv.Itoa(i)+".id", e.Error())
		}

		// if product_qty below the product.min_order_qty
		if v.TotalWeight < v.Product.OrderMinQty {
			o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("order qty"))
		}

		if math.Mod(v.TotalWeight, v.Product.OrderMinQty) != 0 {
			o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("multiple order qty"))
		}

		if v.Product.OrderMaxQty != 0 {
			if v.TotalWeight > v.Product.OrderMaxQty {
				o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("max qty order"))
			}
		}

		//if not exist put in new data
		if box.ID == 0 {
			c.NewProductBox = append(c.NewProductBox, v)
		} else {
			v.Box = box
			existProductBox := &model.BoxItem{Box: v.Box, Status: 1}
			if e := existProductBox.Read("Box", "Status"); e != nil {
				c.ExistProductBox = append(c.ExistProductBox, v)
			} else {
				continue
			}
		}
	}
	return o
}

func (c *createRequestTemplate) Messages() map[string]string {
	messages := map[string]string{
		// "unit_price.required": util.ErrorInputRequired("Unit Price"),
	}

	return messages
}
