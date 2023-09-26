// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createFieldPurchaserRequest : struct to hold create supplier in field purchaser mobile app request data
type createFieldPurchaserRequest struct {
	Code                   string `json:"-"`
	Name                   string `json:"name" valid:"required|alpha_num_space|lte:100"`
	PicName                string `json:"pic_name" valid:"required|alpha_num_space|lte:100"`
	PhoneNumber            string `json:"phone_number" valid:"required|range:8,15|numeric"`
	Address                string `json:"address" valid:"required|lte:350"`
	BlockNumber            string `json:"block_number" valid:"required|lte:10"`
	PaymentMethodID        string `json:"payment_method_id" valid:"required"`
	Rejectable             int8   `json:"rejectable" valid:"required"`
	Returnable             int8   `json:"returnable" valid:"required"`
	Note                   string `json:"note" valid:"alpha_num_space|lte:250"`
	SupplierOrganizationID string `json:"supplier_organization_id" valid:"required"`

	PaymentMethod        *model.PaymentMethod        `json:"-"`
	SupplierOrganization *model.SupplierOrganization `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *createFieldPurchaserRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	c.Code = util.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 3)

	c.PhoneNumber = util.ParsePhoneNumberPrefix(c.PhoneNumber)

	supplier := &model.Supplier{PhoneNumber: c.PhoneNumber}
	if err = supplier.Read("PhoneNumber"); err == nil {
		o.Failure("phone_number.exist", util.ErrorDuplicate("phone number"))
	}

	supplierOrganizationID, err := common.Decrypt(c.SupplierOrganizationID)
	if err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization"))
	}

	c.SupplierOrganization, err = repository.ValidSupplierOrganization(supplierOrganizationID)
	if err != nil {
		o.Failure("supplier_organization_id.invalid", util.ErrorInvalidData("supplier organization"))
	}

	if c.SupplierOrganization.Status != 1 {
		o.Failure("supplier_organization_id.active", util.ErrorActive("supplier organization"))
	}

	paymentMethodID, err := common.Decrypt(c.PaymentMethodID)
	if err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	if c.PaymentMethod, err = repository.ValidPaymentMethod(paymentMethodID); err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createFieldPurchaserRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                     util.ErrorInputRequired("name"),
		"supplier_organization_id.required": util.ErrorSelectRequired("supplier organization"),
		"pic_name.required":                 util.ErrorInputRequired("pic name"),
		"phone_number.required":             util.ErrorInputRequired("phone number"),
		"block_number.required":             util.ErrorInputRequired("block number"),
		"payment_method_id.required":        util.ErrorSelectRequired("payment method"),
		"rejectable.required":               util.ErrorInputRequired("rejectable"),
		"returnable.required":               util.ErrorInputRequired("returnable"),
		"name.alpha_num_space":              util.ErrorAlphaNum("name"),
		"name.lte":                          util.ErrorEqualLess("name", "100"),
		"pic_name.alpha_num_space":          util.ErrorAlphaNum("pic name"),
		"pic_name.lte":                      util.ErrorEqualLess("pic name", "100"),
		"phone_number.numeric":              util.ErrorNumeric("phone number"),
		"phone_number.range":                util.ErrorRangeChar("phone number", "8", "15"),
		"block_number.lte":                  util.ErrorEqualLess("block number", "10"),
		"note.lte":                          util.ErrorEqualLess("note", "250"),
		"address.required":                  util.ErrorInputRequired("address"),
		"address.lte":                       util.ErrorEqualLess("address", "350"),
	}
}
