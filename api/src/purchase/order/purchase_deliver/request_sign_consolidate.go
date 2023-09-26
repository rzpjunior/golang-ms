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

// signConsolidateRequest : struct to hold sign consolidate request data
type signConsolidateRequest struct {
	ConsolidatedPurchaseDeliverID string `json:"consolidated_purchase_deliver_id" valid:"required"`
	Role                          string `json:"role" valid:"required|alpha_num_space|lte:100"`
	Name                          string `json:"name" valid:"required|alpha_num_space|lte:100"`
	Signature                     string `json:"signature" valid:"required"`

	ConsolidatedPurchaseDeliver *model.ConsolidatedPurchaseDeliver `json:"-"`
	Session                     *auth.SessionData                  `json:"-"`
}

// Validate : function to validate sign request data
func (r *signConsolidateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	consolidatedPurchaseDeliverID, err := common.Decrypt(r.ConsolidatedPurchaseDeliverID)
	if err != nil {
		o.Failure("consolidated_purchase_deliver_id.invalid", util.ErrorInvalidData("consolidated purchase deliver id"))
	}

	r.ConsolidatedPurchaseDeliver = &model.ConsolidatedPurchaseDeliver{ID: consolidatedPurchaseDeliverID}
	if err = r.ConsolidatedPurchaseDeliver.Read("ID"); err != nil {
		o.Failure("consolidated_purchase_deliver_id.invalid", util.ErrorInvalidData("consolidated purchase deliver id"))
	}

	// validate if signature is no more than 3 times
	var filter = map[string]interface{}{"consolidated_purchase_deliver_id": consolidatedPurchaseDeliverID}
	var exclude = map[string]interface{}{}

	_, total, err := repository.GetFilterConsolidatedPurchaseDeliverSignature(filter, exclude)
	if err != nil {
		o.Failure("consolidated_purchase_deliver_id.invalid", util.ErrorInvalidData("consolidated purchase deliver id"))
	}

	if total >= 4 {
		o.Failure("consolidated_purchase_deliver_signature.forbidden", util.ErrorEqualLess("signature", "4 times"))
	}

	// validate if logged in user is sourcing admin
	if !(r.Session.Staff.Role.Name == "Sourcing Admin" || r.Session.Staff.Role.Name == "Field Purchaser") {
		o.Failure("role.invalid", util.ErrorRole("add signature", "sourcing admin or field purchaser"))
	}

	// validate if logged in user is the person whom created the surat jalan
	if r.ConsolidatedPurchaseDeliver.CreatedBy.ID != r.Session.Staff.ID {
		o.Failure("user.forbidden", util.ErrorNotValidFor(r.Session.Staff.Name, "consolidated surat jalan"))
	}

	// check if role already sign
	isRoleAlreadySigned, err := repository.IsRoleCPAlreadySigned(consolidatedPurchaseDeliverID, r.Role)
	if err != nil {
		o.Failure("consolidated_purchase_deliver_id.invalid", util.ErrorInvalidData("consolidated purchase deliver"))
	}

	if !isRoleAlreadySigned {
		o.Failure("role.invalid", util.ErrorRoleAlreadySigned(r.Role, r.Name))
	}

	return o
}

func (r *signConsolidateRequest) Messages() map[string]string {
	return map[string]string{
		"consolidated_purchase_deliver_id.required": util.ErrorInputRequired("purchase_deliver_id"),
		"role.required":        util.ErrorInputRequired("role"),
		"name.required":        util.ErrorInputRequired("name"),
		"signature.required":   util.ErrorInputRequired("signature"),
		"role.lte":             util.ErrorEqualLess("role", "100"),
		"name.lte":             util.ErrorEqualLess("name", "100"),
		"role.alpha_num_space": util.ErrorAlphaNum("role"),
		"name.alpha_num_space": util.ErrorAlphaNum("name"),
	}
}
