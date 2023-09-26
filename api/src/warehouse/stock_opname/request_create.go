// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock_opname

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
	CodeStockOpname   string   `json:"-"`
	WarehouseID       string   `json:"warehouse_id" valid:"required"`
	CategoryID        string   `json:"category_id" valid:"required"`
	AreaID            string   `json:"area_id" valid:"required"`
	RecognitionDate   string   `json:"recognition_date" valid:"required"`
	Note              string   `json:"note"`
	Classification    int8     `json:"classification" valid:"required"`
	StockTypeID       int8     `json:"stock_type" valid:"required"`
	StockOpnameItems  []*items `json:"stock_opname_items" valid:"required"`
	RecognitionDateAt time.Time
	Warehouse         *model.Warehouse
	Category          *model.Category
	Session           *auth.SessionData
	StockType         *model.Glossary
}
type items struct {
	ProductID       string  `json:"product_id" valid:"required"`
	FinalStock      float64 `json:"final_stock" valid:"required"`
	AdjustQty       float64
	InitialStock    float64
	OpnameReason    string `json:"opname_reason" valid:"required"`
	OpnameReasonInt int8
	Note            string         `json:"note"`
	Product         *model.Product `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	layout := "2006-01-02"
	var e error
	if c.RecognitionDateAt, e = time.Parse(layout, c.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInvalidData("stock opname date"))
	}
	var filter, exclude map[string]interface{}

	if whID, err := common.Decrypt(c.WarehouseID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		return o
	} else {
		if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		} else {
			if c.Warehouse.Status != int8(1) {
				o.Failure("warehouse_id.invalid", util.ErrorActive("warehouse"))
			}
		}
	}
	if ctID, err := common.Decrypt(c.CategoryID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("category"))
		return o
	} else {
		if c.Category, e = repository.ValidCategory(ctID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("category"))
		} else {
			if c.Category.Status != int8(1) {
				o.Failure("warehouse_id.invalid", util.ErrorActive("category"))
			}
		}
	}

	c.StockType, e = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_int", c.StockTypeID)
	if e != nil {
		o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
		return o
	}

	if len(c.StockOpnameItems) < 0 {
		o.Failure("id.invalid", util.ErrorSelectOne("product"))
	}

	var duplicated = make(map[string]bool)
	var countStockOpname int
	orSelect.Raw("SELECT COUNT(id) FROM stock_opname WHERE category_id = ? AND warehouse_id = ? AND stock_type = ? AND status = 1", c.Category.ID, c.Warehouse.ID, c.StockType.ValueInt).QueryRow(&countStockOpname)
	if countStockOpname > 0 {
		o.Failure("id.invalid", util.ErrorOneActiveSameCategory())
	}
	var classification int8

	if c.Category.GrandParentID == 0 && c.Category.ParentID == 0 {
		classification = 1
	} else if c.Category.GrandParentID != 0 && c.Category.ParentID == 0 {
		classification = 2
	} else {
		classification = 3
	}

	for n, row := range c.StockOpnameItems {
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
					row.Product.Category.Read("ID")
					switch classification {
					case 1:
						if row.Product.Category.GrandParentID != c.Category.ID {
							o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorMustBeSame("product category(c0)", "selected stock opname category"))
						}
					case 2:
						if row.Product.Category.ParentID != c.Category.ID {
							o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorMustBeSame("product category(c1)", "selected stock opname category"))
						}
					case 3:
						if row.Product.Category.ID != c.Category.ID {
							o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorMustBeSame("product category(c2)", "selected stock opname category"))
						}
					default:
						if row.Product.Category.ID != c.Category.ID {
							o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorMustBeSame("product category(c2)", "selected stock opname category"))
						}
					}

					if row.OpnameReason != "" {
						opnameReason, e := repository.GetGlossaryMultipleValue("table", "stock_opname", "attribute", "opname_reason", "value_name", row.OpnameReason)
						if e != nil {
							o.Failure("opname_reason"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("opname_reason"))
						} else {
							row.OpnameReasonInt = opnameReason.ValueInt
						}
					}

					if c.StockType.ValueName == "good stock" {
						orSelect.Raw("SELECT available_stock FROM stock WHERE warehouse_id = ? AND product_id = ? AND status = 1", c.Warehouse.ID, row.Product.ID).QueryRow(&row.InitialStock)
					} else {
						orSelect.Raw("SELECT waste_stock FROM stock WHERE warehouse_id = ? AND product_id = ? AND status = 1", c.Warehouse.ID, row.Product.ID).QueryRow(&row.InitialStock)
					}
				}
				duplicated[row.ProductID] = true
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}
	}
	if len(c.Note) > 250 {
		o.Failure("note.invalid", util.ErrorCharLength("note", 250))
	}
	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"warehouse_id.required":     util.ErrorSelectRequired("warehouse"),
		"category_id.required":      util.ErrorSelectRequired("category"),
		"recognition_date.required": util.ErrorSelectRequired("stock opname date"),
		"area_id.required":          util.ErrorSelectRequired("area"),
		"product_id.required":       util.ErrorSelectRequired("product"),
	}
}
