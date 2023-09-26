// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// signRequest : struct to hold sign request data
type signRequest struct {
	PurchaseOrderID string `json:"purchase_order_id" valid:"required"`
	JobFunction     string `json:"job_function" valid:"required|alpha_num_space|lte:100"`
	Name            string `json:"name" valid:"required|alpha_num_space|lte:100"`
	SignatureURL    string `json:"signature_url" valid:"required"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

// Validate : function to validate sign request data
func (r *signRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	purchaseOrderID, err := common.Decrypt(r.PurchaseOrderID)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
	}

	r.PurchaseOrder, err = repository.ValidPurchaseOrder(purchaseOrderID)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
	}

	// validate if signature is no more than 4 times
	var filter = map[string]interface{}{"purchase_order_id": purchaseOrderID}
	var exclude = map[string]interface{}{}

	_, total, err := repository.GetFilterPurchaseOrderSignature(filter, exclude)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order id"))
	}

	if total >= 4 {
		o.Failure("purchase_order_signature.invalid", util.ErrorEqualLess("signature", "4 times"))
	}

	// check if role already sign
	isRoleAlreadySigned, err := repository.IsRoleAlreadySignedPurchaseOrder(purchaseOrderID, r.JobFunction)
	if err != nil {
		o.Failure("purchase_order_id.invalid", util.ErrorInvalidData("purchase order"))
	}

	if !isRoleAlreadySigned {
		o.Failure("job_function.invalid", util.ErrorRoleAlreadySigned(r.JobFunction, r.Name))
	}

	return o
}

func (r *signRequest) Messages() map[string]string {
	return map[string]string{
		"purchase_order_id.required":   util.ErrorInputRequired("purchase order id"),
		"job_function.required":        util.ErrorInputRequired("job function"),
		"name.required":                util.ErrorInputRequired("name"),
		"signature_url.required":       util.ErrorInputRequired("signature"),
		"job_function.lte":             util.ErrorEqualLess("job function", "100"),
		"name.lte":                     util.ErrorEqualLess("name", "100"),
		"job_function.alpha_num_space": util.ErrorAlphaNum("job function"),
		"name.alpha_num_space":         util.ErrorAlphaNum("name"),
	}
}
