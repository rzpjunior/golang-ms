// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package disposal

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code              string    `json:"-"`
	RecognitionDate   string    `json:"recognition_date" valid:"required"`
	AreaID            string    `json:"area_id" valid:"required"`
	WarehouseID       string    `json:"warehouse_id" valid:"required"`
	Note              string    `json:"note"`
	RecognitionDateAt time.Time `json:"-"`

	Warehouse          *model.Warehouse `json:"-"`
	WasteDisposalItems []*itemRequest   `json:"waste_disposal_items" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type itemRequest struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Uom       string  `json:"uom"`
	Quantity  float64 `json:"dispose_qty"`
	Note      string  `json:"note" valid:"lte:255"`

	Product *model.Product `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var st model.Stock
	var whID int64
	var filter, exclude map[string]interface{}
	var changes bool
	o1 := orm.NewOrm()
	o1.Using("read_only")

	if c.RecognitionDateAt, e = time.Parse("2006-01-02", c.RecognitionDate); e == nil {
		if c.RecognitionDateAt.After(time.Now()) {
			o.Failure("recognition_date.invalid", util.ErrorEqualLess("waste disposal date", "today"))
		}
	} else {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("waste disposal date"))
	}

	if whID, e = common.Decrypt(c.WarehouseID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}
	}

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	var duplicated = make(map[string]bool)

	if e == nil {
		for n, v := range c.WasteDisposalItems {
			var productID int64
			if v.ProductID != "" {
				if !duplicated[v.ProductID] {
					productID, _ = common.Decrypt(v.ProductID)
					if v.Product, e = repository.ValidProduct(productID); e != nil {
						o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
					} else {
						filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.Warehouse.ID, "status": 1}
						if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
							o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
						}
					}
					o1.Raw("select * from stock s where s.product_id = ? and s.warehouse_id = ?", v.Product.ID, c.Warehouse.ID).QueryRow(&st)
					if st.WasteStock < v.Quantity {
						o.Failure("waste_disposal_items.invalid", util.ErrorEqualLess("dispose qty", "waste stock"))
					}
				} else {
					o.Failure("waste_disposal_items.duplicate", util.ErrorDuplicate("product"))
				}
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
			}

			if v.Quantity < 0 {
				o.Failure("waste_disposal_items.invalid", util.ErrorGreater("dispose qty", "0"))
			} else if v.Quantity > 0 {
				changes = true
			}
		}
	}

	//if no quantity was edited at all in the waste disposal
	if changes == false {
		o.Failure("waste_disposal_invalid", util.ErrorNoChanges("waste disposal"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"recognition_date.required": util.ErrorInputRequired("recognition date"),
		"area_id.required":          util.ErrorInputRequired("area"),
		"warehouse_id.required":     util.ErrorInputRequired("warehouse"),
	}

	return messages
}
