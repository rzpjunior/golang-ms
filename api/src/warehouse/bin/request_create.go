// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bin

import (
	"regexp"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold bin create request data
type createRequest struct {
	Code        string  `json:"-"`
	Name        string  `json:"name" valid:"required"`
	WarehouseID string  `json:"warehouse_id" valid:"required"`
	ProductID   string  `json:"product_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Note        string  `json:"note"`

	ServiceTimeStr string `json:"service_time"`
	ServiceTime    int64  `json:"-"`

	Product   *model.Product   `json:"-"`
	Warehouse *model.Warehouse `json:"-"`

	ContainProduct bool  `json:"-"`
	BinAssociated  int64 `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate bin request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	var filter, exclude map[string]interface{}
	var e error

	if c.Code, e = util.CheckTable("bin"); e != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	// warehouse validation
	whID, e := common.Decrypt(c.WarehouseID)
	if e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if c.Warehouse, e = repository.ValidWarehouse(whID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	// service time processing
	serviceTimeArr := strings.Split(c.ServiceTimeStr, ":")
	minutesCheck := regexp.MustCompile(`^[0-9]*$`).MatchString(serviceTimeArr[0])
	if minutesCheck == false {
		o.Failure("minutes.invalid", util.ErrorInvalidData("minutes"))
	}
	secondsCheck := regexp.MustCompile(`^[0-9]*$`).MatchString(serviceTimeArr[1])
	if secondsCheck == false {
		o.Failure("seconds.invalid", util.ErrorInvalidData("seconds"))
	}

	minutes, _ := strconv.Atoi(serviceTimeArr[0])
	seconds, _ := strconv.Atoi(serviceTimeArr[1])
	c.ServiceTime = (int64(minutes) * 60) + int64(seconds)
	if c.ServiceTime <= 0 {
		o.Failure("service_time.invalid", util.ErrorInvalidData("service time"))
	}

	// product id validation
	if len(c.ProductID) > 0 {
		c.ContainProduct = true
		productID, e := common.Decrypt(c.ProductID)
		if e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		if c.Product, e = repository.ValidProduct(productID); e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		}

		//check if product already has associated rack
		stock := &model.Stock{
			Warehouse: c.Warehouse,
			Product:   c.Product,
		}
		if e = orSelect.Read(stock, "Product", "Warehouse"); e != nil {
			o.Failure("stock.invalid", util.ErrorInvalidData("stock"))
		}

		if stock.Bin != nil {
			binID := stock.Bin.ID
			if binID != 0 {
				c.BinAssociated = binID
			}
		}
	}

	// name validation
	filter = map[string]interface{}{"name": c.Name, "warehouse_id": c.Warehouse.ID, "status": 1}
	if _, countBin, err := repository.CheckBinData(filter, exclude); err == nil && countBin != 0 {
		o.Failure("name.invalid", util.ErrorDuplicate("name"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"name.required":         util.ErrorInputRequired("name"),
		"warehouse_id.required": util.ErrorInputRequired("warehouse"),
	}

	return messages
}
