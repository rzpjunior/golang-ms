// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package prospect_supplier

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	Code           string   `json:"-"`
	Name           string   `json:"name" valid:"required"`
	PhoneNumber    string   `json:"phone_number" valid:"required"`
	AltPhoneNumber string   `json:"alt_phone_number"`
	SubDistrictID  string   `json:"sub_district_id" valid:"required"`
	StreetAddress  string   `json:"street_address"`
	Commodity      []string `json:"commodity" valid:"required"`
	PicName        string   `json:"pic_name"`
	PicPhoneNumber string   `json:"pic_phone_number"`
	PicAddress     string   `json:"pic_address"`
	TimeConsent    int8     `json:"time_consent" valid:"required"`
	CommodityStr   string   `json:"-"`

	SubDistrict *model.SubDistrict `json:"-"`

	Session *auth.SessionData `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.Code, err = util.CheckTable("prospect_supplier"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	subDistrictID, _ := common.Decrypt(c.SubDistrictID)
	if c.SubDistrict, err = repository.ValidSubDistrict(subDistrictID); err != nil {
		o.Failure("sub_district_id.invalid", util.ErrorInvalidData("sub district id"))
	}

	c.CommodityStr = strings.Join(c.Commodity, ",")

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":            util.ErrorInputRequired("name"),
		"phone_number.required":    util.ErrorInputRequired("phone number"),
		"sub_district_id.required": util.ErrorInputRequired("sub district"),
		"commodity.required":       util.ErrorInputRequired("commodity"),
		"time_consent.required":    util.ErrorInputRequired("time consent"),
	}
}
