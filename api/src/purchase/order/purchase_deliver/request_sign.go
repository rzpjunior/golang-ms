// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package purchase_deliver

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
	PurchaseDeliverID string `json:"purchase_deliver_id" valid:"required"`
	Role              string `json:"role" valid:"required|alpha_num_space|lte:100"`
	Name              string `json:"name" valid:"required|alpha_num_space|lte:100"`
	Signature         string `json:"signature" valid:"required"`

	PurchaseDeliver *model.PurchaseDeliver `json:"-"`
	Session         *auth.SessionData      `json:"-"`
}

// Validate : function to validate sign request data
func (r *signRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	purchaseDeliverID, err := common.Decrypt(r.PurchaseDeliverID)
	if err != nil {
		o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver id"))
	}

	r.PurchaseDeliver = &model.PurchaseDeliver{ID: purchaseDeliverID}
	if err = r.PurchaseDeliver.Read("ID"); err != nil {
		o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver id"))
	}

	// validate if signature is no more than 3 times
	var filter = map[string]interface{}{"purchase_deliver_id": purchaseDeliverID}
	var exclude = map[string]interface{}{}

	_, total, err := repository.GetFilterPurchaseDeliverSignature(filter, exclude)
	if err != nil {
		o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver id"))
	}

	if total >= 3 {
		o.Failure("purchase_deliver_signature.forbidden", util.ErrorEqualLess("signature", "3 times"))
	}

	// validate if logged in user is sourcing admin
	if !(r.Session.Staff.Role.Name == "Sourcing Admin" || r.Session.Staff.Role.Name == "Field Purchaser") {
		o.Failure("role.invalid", util.ErrorRole("add signature", "sourcing admin or field purchaser"))
	}

	// validate if logged in user is the person whom created the surat jalan
	if r.PurchaseDeliver.CreatedBy != r.Session.Staff.ID {
		o.Failure("user.forbidden", util.ErrorNotValidFor(r.Session.Staff.Name, "surat jalan"))
	}

	// check if role already sign
	isRoleAlreadySigned, err := repository.IsRoleAlreadySigned(purchaseDeliverID, r.Role)
	if err != nil {
		o.Failure("purchase_deliver_id.invalid", util.ErrorInvalidData("purchase deliver"))
	}

	if !isRoleAlreadySigned {
		o.Failure("role.invalid", util.ErrorRoleAlreadySigned(r.Role, r.Name))
	}

	return o
}

func (r *signRequest) Messages() map[string]string {
	return map[string]string{
		"purchase_deliver_id.required": util.ErrorInputRequired("purchase_deliver_id"),
		"role.required":                util.ErrorInputRequired("role"),
		"name.required":                util.ErrorInputRequired("name"),
		"signature.required":           util.ErrorInputRequired("signature"),
		"role.lte":                     util.ErrorEqualLess("role", "100"),
		"name.lte":                     util.ErrorEqualLess("name", "100"),
		"role.alpha_num_space":         util.ErrorAlphaNum("role"),
		"name.alpha_num_space":         util.ErrorAlphaNum("name"),
	}
}
