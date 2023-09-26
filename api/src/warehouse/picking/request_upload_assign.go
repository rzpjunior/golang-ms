// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// uploadAssignRequest : struct to hold picking request data
type uploadAssignRequest struct {
	Code                string    `json:"-"`
	WarehouseID         string    `json:"warehouse_id" valid:"required"`
	RecognitionDate     string    `json:"recognition_date" valid:"required"`
	Note                string    `json:"note"`
	RecognitionDateTime time.Time `json:"-"`

	Warehouse        *model.Warehouse    `json:"-"`
	ItemUploadAssign []*itemUploadAssign `json:"item_upload_assign" valid:"required"`

	DataCorrection map[int]string `json:"-"`

	Session *auth.SessionData `json:"-"`
}

type itemUploadAssign struct {
	ID             string `json:"id"`
	HelperStr      string `json:"helper_str"`
	VendorStr      string `json:"vendor_str"`
	PlanningStr    string `json:"planning_str"`
	SalesOrderCode string `json:"sales_order_code" valid:"required"`

	// data generate for excel
	BusinessTypeStr    string `json:"business_type_str"`
	OrderTypeStr       string `json:"order_type_str"`
	MerchantNameStr    string `json:"merchant_name_str"`
	OrderStatusStr     string `json:"order_status_str"`
	DeliveryCodeStr    string `json:"delivery_code_str"`
	DeliveryStatusStr  string `json:"delivery_status_str"`
	ShippingAddressStr string `json:"shipping_address_str"`
	ProvinceStr        string `json:"province_str"`
	CityStr            string `json:"city_str"`
	DistrictStr        string `json:"district_str"`
	SubDistrictStr     string `json:"sub_district_str"`
	PostalCodeStr      string `json:"postal_code_str"`
	WRTStr             string `json:"wrt_str"`
	OrderWeightStr     string `json:"order_weight_str"`
	PaymentTermStr     string `json:"payment_term_str"`

	SalesOrder    *model.SalesOrder `json:"-"`
	Helper        *model.Staff
	CourierVendor *model.CourierVendor
	// data for picking item entry
	SalesOrderItem []*model.SalesOrderItem `json:"-"`
}

// Validate : function to validate picking request data
func (c *uploadAssignRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	db := orm.NewOrm()
	db.Using("read_only")
	var e error

	if warehouseID, e := common.Decrypt(c.WarehouseID); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	} else {
		if c.Warehouse, e = repository.ValidWarehouse(warehouseID); e != nil {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}
	}

	if c.RecognitionDateTime, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
		o.Failure("recognition_date.invalid", util.ErrorInputRequired("recognition date"))
	}

	c.DataCorrection = make(map[int]string)

	var helperAmount int
	for n, row := range c.ItemUploadAssign {
		if row.HelperStr != "" {
			helperAmount += 1
			//condition for n SO n Helper
			h := &model.Staff{Name: row.HelperStr}

			if e = h.Read("Name"); e != nil {
				c.DataCorrection[n] = util.ErrorInvalidData("helper")
			} else {

				if h.Status != 1 {
					c.DataCorrection[n] = util.ErrorActive("helper")
				}
				if h.Warehouse.ID != c.Warehouse.ID {
					c.DataCorrection[n] = util.ErrorMustBeSame("warehouse of helper", "warehouse of picking order")
				}
				row.Helper = h
			}

			var so *model.SalesOrder
			db.Raw("select * from sales_order so where so.code = ?", row.SalesOrderCode).QueryRow(&so)

			if so == nil {
				c.DataCorrection[n] += " " + util.ErrorInvalidData("sales order")
			} else {
				if row.SalesOrder, e = repository.ValidSalesOrder(so.ID); e != nil {
					if _, ok := c.DataCorrection[n]; ok {
						c.DataCorrection[n] += " " + util.ErrorInvalidData("sales order")
					} else {

						c.DataCorrection[n] = util.ErrorInvalidData("sales order")
					}
				} else {

					var countSO int8
					db.Raw("select count(*) from sales_order so where so.id = ? and so.status in (1,9,12)", row.SalesOrder.ID).QueryRow(&countSO)

					if countSO == 0 {
						if _, ok := c.DataCorrection[n]; ok {
							c.DataCorrection[n] += " " + util.ErrorInvalidData("sales order")
						} else {

							c.DataCorrection[n] = util.ErrorInvalidData("sales order")
						}
					}

					db.Raw("select * from sales_order_item soi where soi.sales_order_id = ?", row.SalesOrder.ID).QueryRows(&row.SalesOrderItem)

					if row.SalesOrder.Warehouse.ID != c.Warehouse.ID {
						if _, ok := c.DataCorrection[n]; ok {
							c.DataCorrection[n] += " " + util.ErrorInvalidData("warehouse of sales order")
						} else {

							c.DataCorrection[n] = util.ErrorInvalidData("warehouse of sales order")
						}
					}
				}
			}
			cv := &model.CourierVendor{Name: row.VendorStr}

			if e = cv.Read("Name"); e != nil {
				c.DataCorrection[n] += util.ErrorInvalidData("vendor")
			} else {
				if cv.Status != 1 {
					c.DataCorrection[n] += util.ErrorActive("vendor")
				}
				row.CourierVendor = cv
			}
		} else {
			if row.VendorStr != "" {
				c.DataCorrection[n] += util.ErrorInvalidData("helper")
			}
		}
	}

	if helperAmount == 0 {
		o.Failure("helper.id", "No Picker Found")
		return o
	}

	return o
}

// Messages : function to return error validation messages
func (c *uploadAssignRequest) Messages() map[string]string {
	messages := map[string]string{
		"warehouse_id.required":       util.ErrorInputRequired("warehouse"),
		"recognition_date.required":   util.ErrorInputRequired("recognition date"),
		"item_upload_assign.required": util.ErrorSalesOrderCannotBeEmpty(),
	}

	return messages
}
