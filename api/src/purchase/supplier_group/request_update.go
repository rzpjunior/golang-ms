// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_group

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	ID                   int64  `json:"-"`
	SupplierCommoditiyID string `json:"supplier_commodity_id" valid:"required"`
	SupplierBadgeID      string `json:"supplier_badge_id" valid:"required"`
	SupplierTypeID       string `json:"supplier_type_id" valid:"required"`

	SupplierCommodity *model.SupplierCommodity `json:"-"`
	SupplierBadge     *model.SupplierBadge     `json:"-"`
	SupplierType      *model.SupplierType      `json:"-"`
	SupplierGroup     *model.SupplierGroup     `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (u *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	u.SupplierGroup, err = repository.GetSupplierGroup("id", u.ID)

	if err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("supplier group"))
	}

	supplierCommodityID, err := common.Decrypt(u.SupplierCommoditiyID)

	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity"))
	}

	u.SupplierCommodity, err = repository.ValidSupplierCommodity(supplierCommodityID)

	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity"))
	}

	if u.SupplierCommodity.Status != 1 {
		o.Failure("supplier_commodity_id.active", util.ErrorActive("supplier commodity"))
	}

	supplierBadgeId, err := common.Decrypt(u.SupplierBadgeID)

	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge"))
	}

	u.SupplierBadge, err = repository.ValidSupplierBadge(supplierBadgeId)

	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge"))
	}

	if u.SupplierBadge.Status != 1 {
		o.Failure("supplier_badge_id.active", util.ErrorActive("supplier badge"))
	}

	supplierTypeID, err := common.Decrypt(u.SupplierTypeID)

	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))

	}

	u.SupplierType, err = repository.ValidSupplierType(supplierTypeID)

	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))
	}

	if u.SupplierType.Status != 1 {
		o.Failure("supplier_type_id.active", util.ErrorActive("supplier type"))
	}

	isDifferent := false

	if u.SupplierGroup.SupplierCommodity != nil {

		if supplierCommodityID != u.SupplierGroup.SupplierCommodity.ID {
			isDifferent = true
		}

	}

	if u.SupplierGroup.SupplierBadge != nil {

		if supplierBadgeId != u.SupplierGroup.SupplierBadge.ID {
			isDifferent = true
		}

	}

	if u.SupplierGroup.SupplierType != nil {

		if supplierTypeID != u.SupplierGroup.SupplierType.ID {
			isDifferent = true
		}

	}

	if isDifferent {
		exist, err := repository.IsExistsSupplierGroup(supplierCommodityID, supplierBadgeId, supplierTypeID)

		if err != nil {
			o.Failure("supplier_group.id", util.ErrorInvalidData("supplier group"))
		}

		if exist {
			o.Failure("supplier_group.id", util.ErrorDuplicate("supplier group"))
		}
	}

	return o
}

func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"supplier_commodity_id.required": util.ErrorSelectRequired("supplier commodity"),
		"supplier_badge_id.required":     util.ErrorSelectRequired("supplier badge"),
		"supplier_type_id.required":      util.ErrorSelectRequired("supplier type"),
	}
}
