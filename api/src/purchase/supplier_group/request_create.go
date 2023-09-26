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

// createRequest : struct to hold supplier relation request data
type createRequest struct {
	SupplierCommoditiyID string `json:"supplier_commodity_id" valid:"required"`
	SupplierBadgeID      string `json:"supplier_badge_id" valid:"required"`
	SupplierTypeID       string `json:"supplier_type_id" valid:"required"`

	SupplierCommodity *model.SupplierCommodity `json:"-"`
	SupplierBadge     *model.SupplierBadge     `json:"-"`
	SupplierType      *model.SupplierType      `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier type request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	supplierCommodityID, err := common.Decrypt(c.SupplierCommoditiyID)

	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity"))
	}

	c.SupplierCommodity, err = repository.ValidSupplierCommodity(supplierCommodityID)

	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity"))
	}

	if c.SupplierCommodity.Status != 1 {
		o.Failure("supplier_commodity_id.active", util.ErrorActive("supplier commodity"))
	}

	supplierBadgeId, err := common.Decrypt(c.SupplierBadgeID)

	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge"))
	}

	c.SupplierBadge, err = repository.ValidSupplierBadge(supplierBadgeId)

	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge"))
	}

	if c.SupplierBadge.Status != 1 {
		o.Failure("supplier_badge_id.active", util.ErrorActive("supplier badge"))
	}

	supplierTypeID, err := common.Decrypt(c.SupplierTypeID)

	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))

	}

	c.SupplierType, err = repository.ValidSupplierType(supplierTypeID)

	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))
	}

	if c.SupplierType.Status != 1 {
		o.Failure("supplier_type_id.active", util.ErrorActive("supplier type"))
	}

	// Check Exist Supplier Group
	exist, err := repository.IsExistsSupplierGroup(supplierCommodityID, supplierBadgeId, supplierTypeID)

	if err != nil {
		o.Failure("supplier_group.id", util.ErrorInvalidData("supplier group"))
	}

	if exist {
		o.Failure("supplier_group.id", util.ErrorDuplicate("supplier group"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"supplier_commodity_id.required": util.ErrorSelectRequired("supplier commodity"),
		"supplier_badge_id.required":     util.ErrorSelectRequired("supplier badge"),
		"supplier_type_id.required":      util.ErrorSelectRequired("supplier type"),
	}
}
