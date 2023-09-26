// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package supplier_organization

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type createRequest struct {
	Code                string `json:"-"`
	SupplierCommodityID string `json:"supplier_commodity_id" valid:"required"`
	SupplierBadgeID     string `json:"supplier_badge_id" valid:"required"`
	SupplierTypeID      string `json:"supplier_type_id" valid:"required"`
	PurchaseTermID      string `json:"purchase_term_id" valid:"required"`
	SubDistrictID       string `json:"sub_district_id" valid:"required"`
	Name                string `json:"name" valid:"required|alpha_num_space|lte:100"`
	Address             string `json:"address" valid:"required|lte:350"`
	Note                string `json:"note" valid:"lte:250"`

	SupplierCommodity *model.SupplierCommodity `json:"-"`
	SupplierBadge     *model.SupplierBadge     `json:"-"`
	SupplierType      *model.SupplierType      `json:"-"`
	PurchaseTerm      *model.PurchaseTerm      `json:"-"`
	SubDistrict       *model.SubDistrict       `json:"-"`
	Session           *auth.SessionData        `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("supplier_organization"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	supplierCommodityID, err := common.Decrypt(c.SupplierCommodityID)
	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity id"))
	}

	c.SupplierCommodity, err = repository.ValidSupplierCommodity(supplierCommodityID)
	if err != nil {
		o.Failure("supplier_commodity_id.invalid", util.ErrorInvalidData("supplier commodity id"))
	}

	supplierBadgeID, err := common.Decrypt(c.SupplierBadgeID)
	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge id"))
	}

	c.SupplierBadge, err = repository.ValidSupplierBadge(supplierBadgeID)
	if err != nil {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge id"))
	}

	supplierTypeID, err := common.Decrypt(c.SupplierTypeID)
	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type id"))
	}

	c.SupplierType, err = repository.ValidSupplierType(supplierTypeID)
	if err != nil {
		o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type id"))
	}

	purchaseTermID, err := common.Decrypt(c.PurchaseTermID)
	if err != nil {
		o.Failure("purchase_term_id.invalid", util.ErrorInvalidData("purchase term id"))
	}

	c.PurchaseTerm, err = repository.ValidPurchaseTerm(purchaseTermID)
	if err != nil {
		o.Failure("purchase_term_id.invalid", util.ErrorInvalidData("purchase term id"))
	}

	subDistrictID, err := common.Decrypt(c.SubDistrictID)
	if err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district id"))
	}

	c.SubDistrict, err = repository.ValidSubDistrict(subDistrictID)
	if err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district id"))
	}

	exists, err := repository.IsExistsSupplierOrganization(c.Name)

	if err != nil {
		o.Failure("name", util.ErrorInvalidData("supplier organization"))
	}

	if exists {
		o.Failure("name", util.ErrorDuplicate("supplier organization"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                  util.ErrorInputRequired("name"),
		"name.alpha_num_space":           util.ErrorAlphaNum("name"),
		"name.lte":                       util.ErrorEqualLess("name", "100"),
		"address.required":               util.ErrorInputRequired("address"),
		"address.lte":                    util.ErrorEqualLess("address", "350"),
		"supplier_commodity_id.required": util.ErrorSelectRequired("supplier commodity"),
		"supplier_badge_id.required":     util.ErrorSelectRequired("supplier badge"),
		"supplier_type_id.required":      util.ErrorSelectRequired("supplier type"),
		"purchase_term_id.required":      util.ErrorSelectRequired("payment term"),
		"sub_district_id.required":       util.ErrorSelectRequired("sub district"),
		"note.lte":                       util.ErrorEqualLess("note", "250"),
	}
}
