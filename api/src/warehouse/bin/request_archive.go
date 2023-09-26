// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package bin

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// archieveRequest : struct to routing set request data
type archiveRequest struct {
	ID           int64 `json:"-" valid:"required"`
	ContainStock bool  `json:"-"`

	Session *auth.SessionData `json:"-"`
	Stock   *model.Stock      `json:"-"`
}

// Validate : function to validate routing request data
func (a *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	bin, err := repository.ValidBin(a.ID)
	if err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("id"))
	}

	if bin.Status != 1 {
		o.Failure("status.inactive", util.ErrorActive("status"))
	}

	err = bin.Warehouse.Read("ID")
	if err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if bin.Product.ID != 0 {
		stock := &model.Stock{
			Warehouse: bin.Warehouse,
			Bin:       bin,
		}
		if err = orSelect.Read(stock, "Warehouse", "bin"); err != nil {
			o.Failure("stock.invalid", util.ErrorInvalidData("stock"))
		}
		if stock.Product.ID != 0 {
			a.ContainStock = true
			a.Stock = stock
		}
	}

	return o
}

// Messages : function to return error validation messages
func (u *archiveRequest) Messages() map[string]string {
	messages := map[string]string{
		"bin_id.required": util.ErrorInputRequired("bin id"),
	}

	return messages
}
