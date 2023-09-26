// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// groupingSalesOrderRequest : struct to hold picking assign request data
type groupingSalesOrderRequest struct {
	ProductID     string `json:"product_id" valid:"required"`
	PickingListID string `json:"picking_list_id" valid:"required"`

	Product     *model.Product     `json:"-"`
	PickingList *model.PickingList `json:"-"`

	PickingListFinal map[string]ListPl

	PickingOrder       *model.PickingOrder `json:"-"`
	PackRecommendation []*model.PackRecommendation

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign request data
func (r *groupingSalesOrderRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error

	productID, _ := common.Decrypt(r.ProductID)
	pickingListID, _ := common.Decrypt(r.PickingListID)

	if r.Product, err = repository.ValidProduct(productID); err != nil {
		o.Failure("product.invalid", util.ErrorInvalidData("product"))
	}

	if r.PickingList, err = repository.ValidPickingList(pickingListID); err != nil {
		o.Failure("picking_list.invalid", util.ErrorInvalidData("picking list"))
	}

	return o
}

// Messages : function to return error validation messages
func (r *groupingSalesOrderRequest) Messages() map[string]string {
	messages := map[string]string{
		"product_id.required":      util.ErrorInputRequired("product"),
		"picking_list_id.required": util.ErrorInputRequired("picking list"),
	}

	return messages
}
