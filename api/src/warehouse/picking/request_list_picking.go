// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// listRequestPicking : struct to hold picking assign request data
type listRequestPicking struct {
	WarehouseID         string    `json:"warehouse_id" valid:"required"`
	HelperID            string    `json:"helper_id" valid:"required"`
	DeliveryDate        string    `json:"delivery_date" valid:"required"`
	Query               string    `json:"query"`
	FilterPickingList   string    `json:"filter_picking_list"`
	FilterPickingStatus string    `json:"filter_picking_status"`
	RecognitionDateTime time.Time `json:"-"`

	Staff     *model.Staff     `json:"-"`
	Warehouse *model.Warehouse `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *listRequestPicking) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	wID, _ := common.Decrypt(r.WarehouseID)
	hID, _ := common.Decrypt(r.HelperID)

	if r.Warehouse, err = repository.ValidWarehouse(wID); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
	}

	if r.Staff, err = repository.ValidStaff(hID); err != nil {
		o.Failure("helper.invalid", util.ErrorInvalidData("helper"))
	}

	return o
}

// Messages : function to return error validation messages
func (r *listRequestPicking) Messages() map[string]string {
	messages := map[string]string{
		"delivery_date.required": util.ErrorInputRequired("delivery date"),
		"warehouse_id.required":  util.ErrorInputRequired("warehouse"),
		"helper_id.required":     util.ErrorInputRequired("helper"),
	}

	return messages
}
