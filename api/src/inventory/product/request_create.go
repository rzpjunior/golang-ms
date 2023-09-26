// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold product request data
type createRequest struct {
	Code            string   `json:"code" valid:"required"`
	Name            string   `json:"name" valid:"required"`
	Note            string   `json:"note"`
	Description     string   `json:"description"`
	Weight          float64  `json:"weight" valid:"required|gt:0"`
	WarehouseSto    []string `json:"warehouse_sto" valid:"required"`
	WarehousePur    []string `json:"warehouse_pur"`
	TagProduct      []string `json:"product_tag"`
	Storability     int8     `json:"storability" valid:"required"`
	Purchasability  int8     `json:"purchasability" valid:"required"`
	UomId           string   `json:"uom_id" valid:"required"`
	CategoryId      string   `json:"category_id" valid:"required"`
	Images          []string `json:"images" valid:"required"`
	UnivProductCode string   `json:"up_code"`
	OrderMinQty     float64  `json:"order_min_qty" valid:"required|gt:0"`
	SparePercentage float64  `json:"spare_percentage"`
	Taxable         int8     `json:"taxable" valid:"required"`
	TaxPercentage   float64  `json:"tax_percentage"`
	FragileGoods    int8     `json:"fragile_goods"`

	WarehouseStoStr string         `json:"-"`
	WarehousePurStr string         `json:"-"`
	WrhStoExist     map[int64]int8 `json:"-"`
	WrhPurExist     map[int64]int8 `json:"-"`

	Uom       *model.Uom         `json:"-"`
	Category  *model.Category    `json:"-"`
	Warehouse []*model.Warehouse `json:"-"`
	PriceSet  []*model.PriceSet  `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate product request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error
	var filter, exclude map[string]interface{}
	wrhStoExist := make(map[int64]int8)
	wrhPurExist := make(map[int64]int8)

	if dbredis.Redis.CheckExistKey("warehouse_creation") {
		o.Failure("warehouse_creation.invalid", util.ErrorCreationInProgress("warehouse"))
	}

	uomId, e := common.Decrypt(c.UomId)
	if e != nil {
		o.Failure("uom.invalid", util.ErrorInvalidData("uom"))
	}
	if c.Uom, e = repository.ValidUom(uomId); e != nil {
		o.Failure("uom.invalid", util.ErrorInvalidData("uom"))
	}
	if c.Uom.Status != int8(1) {
		o.Failure("uom.active", util.ErrorActive("uom"))
	}

	categoryId, e := common.Decrypt(c.CategoryId)
	if e != nil {
		o.Failure("category.invalid", util.ErrorInvalidData("category"))
	}
	if c.Category, e = repository.ValidCategory(categoryId); e != nil {
		o.Failure("category.invalid", util.ErrorInvalidData("category"))
	}
	if c.Category.Status != int8(1) {
		o.Failure("category.active", util.ErrorActive("category"))
	}

	filter = map[string]interface{}{"name": c.Name, "uom_id": uomId}
	exclude = map[string]interface{}{"status": int8(3)}
	if _, countName, err := repository.CheckProductData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name.unique", util.ErrorUniqueProduct())
	}

	filterCode := map[string]interface{}{"code": c.Code}
	excludeStatus := map[string]interface{}{"status": int8(3)}
	if _, countCode, err := repository.CheckProductData(filterCode, excludeStatus); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	} else if countCode > 0 {
		o.Failure("code", util.ErrorDuplicate("code"))
	}

	if len(c.WarehouseSto) == 0 {
		o.Failure("warehouse_sto.invalid", util.ErrorSelectOne("warehouse sto"))
	}

	if len(strings.TrimSpace(c.Code)) != 8 {
		o.Failure("code.invalid", util.ErrorCharLength("code", 8))
	}

	if c.Purchasability == 1 && len(c.WarehousePur) == 0 {
		o.Failure("warehouse_pur.invalid", util.ErrorSelectOne("warehouse pur"))
	}

	for i, v := range c.WarehouseSto {
		v = common.Encrypt(v)
		c.WarehouseStoStr = c.WarehouseStoStr + v + ","
		c.WarehouseSto[i] = v

		warehouseID, _ := strconv.Atoi(v)
		wrhStoExist[int64(warehouseID)] = 1
	}
	c.WarehouseStoStr = strings.TrimSuffix(c.WarehouseStoStr, ",")
	c.WrhStoExist = wrhStoExist

	for i, v := range c.WarehousePur {
		v = common.Encrypt(v)
		c.WarehousePurStr = c.WarehousePurStr + v + ","
		warehouseID, _ := strconv.Atoi(v)
		if _, isExist := c.WrhStoExist[int64(warehouseID)]; !isExist {
			o.Failure("warehouse_sto.invalid", util.ErrorMustContain("warehouse storability", "warehouse purchasability"))
		}
		c.WarehousePur[i] = v
		wrhPurExist[int64(warehouseID)] = 1
	}
	c.WarehousePurStr = strings.TrimSuffix(c.WarehousePurStr, ",")
	c.WrhPurExist = wrhPurExist

	filter = map[string]interface{}{"status": 1}
	exclude = map[string]interface{}{}
	if c.Warehouse, _, err = repository.CheckWarehousesData(filter, exclude); err != nil {
		o.Failure("status.invalid", util.ErrorInvalidData("status"))
	}

	if c.PriceSet, err = repository.GetAllPriceSets(); err != nil {
		o.Failure("status.invalid", "Failed to Get All Price Sets")
	}
	glossary := &model.Glossary{
		Table:     "product",
		Attribute: "fragile_goods",
		ValueName: "no",
	}
	err = glossary.Read("table", "attribute", "value_name")
	no_fragile_goods := glossary.ValueInt

	if c.FragileGoods != 0 {
		glossary := &model.Glossary{
			Table:     "product",
			Attribute: "fragile_goods",
			ValueInt:  c.FragileGoods,
		}
		err = glossary.Read("table", "attribute", "value_int")
		if err != nil {
			o.Failure("fragile_goods.invalid", util.ErrorInvalidData("fragile goods"))
		}
	} else {
		c.FragileGoods = no_fragile_goods
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":           util.ErrorInputRequired("product name"),
		"weight.required":         util.ErrorInputRequired("weight"),
		"tag_product.required":    util.ErrorInputRequired("tag product"),
		"uom_id.required":         util.ErrorInputRequired("uom"),
		"category_id.required":    util.ErrorInputRequired("category"),
		"storability.required":    util.ErrorInputRequired("storability"),
		"purchasability.required": util.ErrorInputRequired("purchasability"),
		"warehouse_sto.required":  util.ErrorInputRequired("storability assigned warehouse"),
		"images.required":         util.ErrorInputRequired("images"),
		"weight.gt":               util.ErrorGreater("weight", "0"),
		"order_min_qty.gt":        util.ErrorGreater("minimum order quantity", "0"),
	}
}
