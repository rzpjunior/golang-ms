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

type updateRequest struct {
	ID                     int64  `json:"-" valid:"required"`
	Name                   string `json:"name" valid:"required"`
	PicName                string `json:"pic_name" valid:"required"`
	PhoneNumber            string `json:"phone_number" valid:"required|numeric|range:8,15"`
	AltPhoneNumber         string `json:"alt_phone_number" valid:"numeric"`
	Email                  string `json:"email" valid:"email"`
	Address                string `json:"address" valid:"required"`
	TermPaymentPurID       string `json:"term_payment_pur_id" valid:"required"`
	PaymentMethodID        string `json:"payment_method_id" valid:"required"`
	SubDistrictID          string `json:"sub_district_id" valid:"required"`
	Note                   string `json:"note"`
	SupplierTypeID         string `json:"supplier_type_id" valid:"required"`
	SupplierBadgeID        string `json:"supplier_badge_id" valid:"required"`
	SupplierCommodityID    string `json:"supplier_commodity_id" valid:"required"`
	SupplierOrganizationID string `json:"supplier_organization_id"`
	BlockNumber            string `json:"block_number" valid:"lte:10"`
	Rejectable             int8   `json:"rejectable" valid:"required"`
	Returnable             int8   `json:"returnable" valid:"required"`

	TermPaymentPur       *model.PurchaseTerm         `json:"-"`
	PaymentMethod        *model.PaymentMethod        `json:"-"`
	Supplier             *model.Supplier             `json:"-"`
	SubDistrict          *model.SubDistrict          `json:"-"`
	Area                 *model.Area                 `json:"-"`
	SupplierType         *model.SupplierType         `json:"-"`
	SupplierBadge        *model.SupplierBadge        `json:"-"`
	SupplierCommodity    *model.SupplierCommodity    `json:"-"`
	SupplierOrganization *model.SupplierOrganization `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	c.Supplier = &model.Supplier{Name: c.Name}
	if err = c.Supplier.Read("Name"); err == nil && c.Supplier.ID != c.ID {
		o.Failure("name", util.ErrorDuplicate("name"))
	}

	supplierCommodityID, err := common.Decrypt(c.SupplierCommodityID)

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

	if c.SupplierType.ID == 7 {
		if c.BlockNumber == "" {
			o.Failure("block_number.required", util.ErrorInputRequired("block number"))
		}
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number", util.ErrorCharLength("phone number", 8))
	}

	if len(c.AltPhoneNumber) > 0 && len(c.AltPhoneNumber) < 8 {
		o.Failure("alt_phone_number", util.ErrorCharLength("alternative phone number", 8))
	}

	if c.PhoneNumber == c.AltPhoneNumber {
		o.Failure("alt_phone_number", util.ErrorInputCannotBeSame("phone number", "alternative phone number"))
	}

	termPaymentPurID, err := common.Decrypt(c.TermPaymentPurID)

	if err != nil {
		o.Failure("term_payment_pur_id.invalid", util.ErrorInvalidData("purchase term"))
	}

	if c.TermPaymentPur, err = repository.ValidPurchaseTerm(termPaymentPurID); err != nil {
		o.Failure("term_payment_pur_id.invalid", util.ErrorInvalidData("purchase term"))
	}

	paymentMethodID, err := common.Decrypt(c.PaymentMethodID)

	if err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	if c.PaymentMethod, err = repository.ValidPaymentMethod(paymentMethodID); err != nil {
		o.Failure("payment_method_id.invalid", util.ErrorInvalidData("payment method"))
	}

	subDistrictId, err := common.Decrypt(c.SubDistrictID)

	if err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	}

	if c.SubDistrict, err = repository.ValidSubDistrict(subDistrictId); err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district"))
	}

	if c.SubDistrict.Status != int8(1) {
		o.Failure("sub_district_id.active", util.ErrorActive("sub district"))
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

	if supplierBadgeId == 0 {
		o.Failure("supplier_badge_id.invalid", util.ErrorInvalidData("supplier badge"))
	}

	if supplierBadgeId != 0 && supplierCommodityID != 0 && supplierTypeID != 0 {
		filterGroup := map[string]interface{}{"supplier_badge_id": c.SupplierBadge.ID, "supplier_commodity_id": c.SupplierCommodity.ID, "supplier_type_id": c.SupplierType.ID}
		exclude := map[string]interface{}{}

		// Valid Grouping Supplier Commodity - Supplier Badge
		_, total, err := repository.CheckValidSupplierGroup(filterGroup, exclude)

		if err != nil {
			o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))
		}

		if total == 0 {
			o.Failure("supplier_type_id.invalid", util.ErrorInvalidData("supplier type"))
		}
	}

	if c.SupplierOrganizationID != "" {

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
	}

	if c.AltPhoneNumber != "" {
		if len(c.AltPhoneNumber) > 15 || len(c.AltPhoneNumber) < 8 {
			o.Failure("alt_phone_number.range", util.ErrorRangeChar("alt_phone_number", "8", "15"))
		}
	}

	c.PhoneNumber = util.ParsePhoneNumberPrefix(c.PhoneNumber)
	c.AltPhoneNumber = util.ParsePhoneNumberPrefix(c.AltPhoneNumber)

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                util.ErrorInputRequired("name"),
		"pic_name.required":            util.ErrorInputRequired("pic name"),
		"phone_number.required":        util.ErrorInputRequired("phone number"),
		"address.required":             util.ErrorInputRequired("address"),
		"term_payment_pur_id.required": util.ErrorInputRequired("purchase payment"),
		"payment_method_id.required":   util.ErrorInputRequired("payment method"),
		"sub_district_id.required":     util.ErrorInputRequired("sub district"),
		"supplier_type_id.required":    util.ErrorSelectRequired("supplier type"),
		"supplier_badge_id.required":   util.ErrorSelectRequired("supplier badge"),
		"rejectable.required":          util.ErrorInputRequired("rejectable"),
		"returnable.required":          util.ErrorInputRequired("returnable"),
		"block_number.lte":             util.ErrorEqualLess("block number", "10"),
		"phone_number.numeric":         util.ErrorNumeric("phone number"),
		"alt_phone_number.numeric":     util.ErrorNumeric("alt phone number"),
		"phone_number.range":           util.ErrorRangeChar("phone number", "8", "15"),
	}
}
