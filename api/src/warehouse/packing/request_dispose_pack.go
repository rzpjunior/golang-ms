// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// DisposePackRequest: struct generate picking
type DisposePackRequest struct {
	ID        int64   `json:"-"`
	ProductID string  `json:"product_id" valid:"required"`
	PackType  float64 `json:"pack_type" valid:"required"`
	Dispose   int     `json:"dispose"`

	PackingOrder *model.PackingOrder `json:"-"`
	Product      *model.Product      `json:"-"`

	ResponseData *ResponseData     `json:"-"`
	Session      *auth.SessionData `json:"-"`
}

// Validate : function to validate
func (c *DisposePackRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var err error
	var pID int64

	if c.PackingOrder, err = repository.GetPackingOrder("ID", c.ID); err != nil {
		o.Failure("id_invalid", util.ErrorInvalidData("packing order"))
		return o
	}
	// region product definition
	if pID, err = common.Decrypt(c.ProductID); err != nil {
		o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		return o
	}

	if c.Product, err = repository.ValidProduct(pID); err != nil {
		o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
		return o
	}
	if err = c.Product.Uom.Read("ID"); err != nil {
		o.Failure("uom_id.invalid", util.ErrorInvalidData("uom"))
		return o
	}

	o1.LoadRelated(c.Product, "ProductImage", 1)
	// endregion

	return o
}

// Messages : function to return error validation messages
func (c *DisposePackRequest) Messages() map[string]string {
	messages := map[string]string{
		"product_id.required": util.ErrorInputRequired("product"),
		"pack_type.required":  util.ErrorInputRequired("pack type"),
	}

	return messages
}
