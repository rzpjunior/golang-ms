// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// assignPickerRequest : struct to hold picking assign request data
type assignPickerRequest struct {
	ID              int64     `json:"-"`
	PickerID        []string  `json:"picker_id"`
	PickerCapacity  int64     `json:"picker_max_weight"`
	AssignTimeStamp time.Time `json:"-"`

	PickerString          string           `json:"-"`
	ErrorsSalesOrderItems []int64          `json:"-"`
	Picker                *model.Staff     `json:"-"`
	Warehouse             *model.Warehouse `json:"-"`

	PickingList *model.PickingList `json:"-"`
	Session     *auth.SessionData  `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *assignPickerRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var (
		leadPicker         int64
		err                error
		filter, exclude    map[string]interface{}
		weight             float64
		maxWeight          float64
		pickingOrderAssign []*model.PickingOrderAssign
		pickingOrderItem   []*model.PickingOrderItem
	)

	if r.PickingList, err = repository.ValidPickingList(r.ID); err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}
	if r.PickingList.Status != 1 {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}

	o1.Raw("SELECT staff_id FROM picking_order_assign WHERE picking_list_id = ? limit 1", r.PickingList.ID).QueryRow(&leadPicker)
	if leadPicker == 0 {
		o.Failure("lead_picker_id.invalid", util.ErrorInvalidData("lead picker"))
		return o
	}

	// routing type
	if len(r.PickerID) > 0 {
		for i, v := range r.PickerID {
			pickerID, err := common.Decrypt(v)
			if err != nil {
				o.Failure("picker_id"+strconv.Itoa(i), util.ErrorInvalidData("picker"))
			}

			if r.Picker, err = repository.GetStaff("id", pickerID); err != nil {
				o.Failure("picker_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("picker"))
			}

			if r.Picker.Status != 1 {
				o.Failure("picker_id"+strconv.Itoa(i)+".invalid", util.ErrorActiveInd("picker"))
			}

			pickerString := strconv.Itoa(int(pickerID))
			r.PickerString += pickerString + ","
		}
		r.PickerString = strings.TrimSuffix(r.PickerString, ",")

		if r.PickerCapacity == 0 {
			o.Failure("picker_capacity.invalid", util.ErrorGreater("picker capacity", "0"))
		}

		// warehouse bin info is required
		r.Warehouse = &model.Warehouse{ID: r.PickingList.Warehouse.ID}
		if err = r.Warehouse.Read("id"); err != nil {
			o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
		}
		if r.Warehouse.BinInfo == nil {
			o.Failure("bin_info.invalid", util.ErrorInvalidData("bin info"))
			return o
		}
		if err = r.Warehouse.BinInfo.Read("ID"); err != nil {
			o.Failure("bin_info.invalid", util.ErrorInvalidData("bin info"))
			return o
		}

		filter = map[string]interface{}{"picking_list_id": r.PickingList.ID}
		pickingOrderAssign, _, err = repository.CheckPickingOrderAssignData(filter, exclude)
		if err != nil {
			o.Failure("picking_order_assign.invalid", util.ErrorNotFound("picking order assign"))
		}

		for _, v1 := range pickingOrderAssign {
			filter = map[string]interface{}{"picking_order_assign_id": v1.ID}
			pickingOrderItem, _, err = repository.CheckPickingOrderItemData(filter, exclude)
			if err != nil {
				o.Failure("picking_order_item.invalid", util.ErrorNotFound("picking order item"))
			}

			for _, v2 := range pickingOrderItem {
				salesOrderItem := &model.SalesOrderItem{SalesOrder: v1.SalesOrder, Product: v2.Product}
				if err = salesOrderItem.Read("salesorder", "product"); err != nil {
					weight = v2.OrderQuantity
				} else {
					weight = math.Ceil(salesOrderItem.Weight)
				}
				if weight > maxWeight {
					maxWeight = weight
				}

				product := &model.Product{ID: v2.Product.ID}
				err = product.Read("ID")
				if err != nil {
					o.Failure("product.not_found", util.ErrorNotFound("product"))
				}

				stock := &model.Stock{
					Product:   product,
					Warehouse: r.Warehouse,
				}
				if err = o1.Read(stock, "Product", "Warehouse"); err != nil {
					o.Failure("stock.invalid", util.ErrorInvalidData("stock"))
				} else {
					if stock.Bin == nil {
						r.ErrorsSalesOrderItems = append(r.ErrorsSalesOrderItems, v2.Product.ID)
					} else {
						if stock.Bin.ID == 0 {
							r.ErrorsSalesOrderItems = append(r.ErrorsSalesOrderItems, v2.Product.ID)
						} else {
							bin := &model.Bin{
								ID:        stock.Bin.ID,
								Warehouse: r.Warehouse,
								Product:   product,
							}
							if err = bin.Read("ID", "Warehouse", "Product"); err != nil {
								r.ErrorsSalesOrderItems = append(r.ErrorsSalesOrderItems, v2.Product.ID)
							}
						}
					}
				}
			}
		}

		// if there's a product that doesn't have bin associated return all the products
		if len(r.ErrorsSalesOrderItems) > 0 {
			var errorSalesOrderItemsString string
			for _, v := range r.ErrorsSalesOrderItems {
				product := &model.Product{ID: v}
				if err = product.Read("ID"); err != nil {
					o.Failure("product.not_found", util.ErrorNotFound("product"))
				}

				errorSalesOrderItemsString += "Product " + product.Name + " Doesn't have bin associated. | "
			}

			errorSalesOrderItemsString = errorSalesOrderItemsString[:len(errorSalesOrderItemsString)-3]
			o.Failure("no_bin", errorSalesOrderItemsString)
		}

		if maxWeight > float64(r.PickerCapacity) {
			maxWeightString := fmt.Sprintf("%.0f", maxWeight)
			o.Failure("picker_max_weight.invalid", util.ErrorEqualGreater("picker max weight", maxWeightString))
		}
		r.AssignTimeStamp = time.Time{}
	} else {
		r.PickerCapacity = 0
	}

	return o
}

// Messages : function to return error validation messages
func (r *assignPickerRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
