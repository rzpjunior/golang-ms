// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entry

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type createRequest struct {
	CodeWasteEntry    string   `json:"-"`
	WarehouseID       string   `json:"warehouse_id" valid:"required"`
	AreaID            string   `json:"area_id" valid:"required"`
	RecognitionDate   string   `json:"recognition_date" valid:"required"`
	Note              string   `json:"note"`
	WasteEntryItems   []*items `json:"waste_entry_items" valid:"required"`
	RecognitionDateAt time.Time

	Warehouse *model.Warehouse
	Session   *auth.SessionData
}
type items struct {
	ID             string  `json:"id"`
	ProductID      string  `json:"product_id" valid:"required"`
	WasteStock     float64 `json:"waste_stock" valid:"required"`
	InitialStock   float64
	WasteReason    string                `json:"waste_reason"`
	Note           string                `json:"note"`
	Product        *model.Product        `json:"-"`
	WasteEntryItem *model.WasteEntryItem `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	var changes bool
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if c.RecognitionDateAt, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("waste entry date"))
	}

	if whID, err := common.Decrypt(c.WarehouseID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		} else {
			if c.Warehouse.Status != int8(1) {
				o.Failure("warehouse_id.invalid", util.ErrorActive("warehouse"))
			}
		}
	}

	// Get today date
	today := time.Now()

	if c.RecognitionDateAt.After(today) {
		o.Failure("recognition_date.invalid", util.ErrorEqualLess("waste entry date", "today date"))
	}

	if len(c.WasteEntryItems) < 0 {
		o.Failure("id.invalid", util.ErrorSelectOne("product"))
	}

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	var duplicated = make(map[string]bool)

	for n, row := range c.WasteEntryItems {
		var productID int64
		if row.ProductID != "" {
			if !duplicated[row.ProductID] {
				productID, _ = common.Decrypt(row.ProductID)
				if row.Product, e = repository.ValidProduct(productID); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
				} else {
					filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.Warehouse.ID, "status": 1}
					if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
						o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
					}

					if row.WasteStock < 0 {
						o.Failure("qty"+strconv.Itoa(n)+".greater", util.ErrorGreater("product quantity", "0"))
					} else if row.WasteStock > 0 {
						// blocked reason is empty
						if row.WasteReason == "" {
							o.Failure("waste_reason"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("waste_reason"))
						} else {
							// check reason exiting in glossary
							_, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "waste_reason", "value_name", row.WasteReason)
							if e != nil {
								o.Failure("waste_reason"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("waste_reason"))
							}
						}

						changes = true
					}

					orSelect.Raw("SELECT available_stock FROM stock WHERE warehouse_id = ? AND product_id = ? AND status = 1", c.Warehouse.ID, row.Product.ID).QueryRow(&row.InitialStock)
				}
				duplicated[row.ProductID] = true

			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}
	}

	//if no quantity was edited at all in the waste entry
	if changes == false {
		o.Failure("waste_entry_invalid", util.ErrorNoChanges("waste entry"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"warehouse_id.required":     util.ErrorSelectRequired("warehouse"),
		"area_id.required":          util.ErrorInputRequired("area"),
		"recognition_date.required": util.ErrorSelectRequired("waste entry date"),
		"product_id.required":       util.ErrorSelectRequired("product"),
	}
}
