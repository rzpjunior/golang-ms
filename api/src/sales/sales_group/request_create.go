// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_group

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
	"strings"
)

// createRequest : struct to hold sales group request data
type createRequest struct {
	Code           string   `json:"-"`
	Name           string   `json:"name" valid:"required"`
	BusinessTypeID string   `json:"business_type_id" valid:"required"`
	SlsManID       string   `json:"sls_man_id" valid:"required"`
	AreaID         string   `json:"area_id" valid:"required"`
	SubDistrict    []string `json:"sub_district_id" valid:"required"`
	SubDistrictStr string   `json:"-"`
	CityStr        string   `json:"-"`

	BusinessType *model.BusinessType `json:"-"`
	SlsMan       *model.Staff        `json:"-"`
	Area         *model.Area         `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate sales group request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var filter, exclude map[string]interface{}
	var businessTypeID, slsManID, areaID int64
	var subDistrictExist = make(map[string]int8)
	var duplicated = make(map[int64]bool)
	var cityList []string

	if r.Code, err = util.CheckTable("sales_group"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if len(r.Name) > 50 {
		o.Failure("name", util.ErrorCharLength("name", 50))
	}

	filter = map[string]interface{}{"name": r.Name, "status": int8(1)}
	exclude = map[string]interface{}{"status": int8(2)}
	if _, countName, err := repository.CheckSalesGroupData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name.unique", util.ErrorUnique("name"))
	}

	if businessTypeID, err = common.Decrypt(r.BusinessTypeID); err != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	if r.BusinessType, err = repository.ValidBusinessType(businessTypeID); err != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	if r.BusinessType.Status != 1 {
		o.Failure("business_type_id.inactive", util.ErrorActive("business type"))
	}

	if slsManID, err = common.Decrypt(r.SlsManID); err != nil {
		o.Failure("sls_man_id.invalid", util.ErrorInvalidData("sales manager"))
	}

	if r.SlsMan, err = repository.ValidStaff(slsManID); err != nil {
		o.Failure("sls_man_id.invalid", util.ErrorInvalidData("sales manager"))
	}

	if areaID, err = common.Decrypt(r.AreaID); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	}

	if r.Area, err = repository.ValidArea(areaID); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	}

	for i, v := range r.SubDistrict {
		v = common.Encrypt(v)

		idSubDist, _ := strconv.ParseInt(v, 10, 64)

		// get district
		subDist := &model.SubDistrict{ID: idSubDist}
		subDist.Read("ID")

		// get city
		dist := &model.District{ID: subDist.District.ID}
		dist.Read("ID")

		if !duplicated[dist.City.ID] {
			cityList = append(cityList, strconv.FormatInt(dist.City.ID, 10))
			duplicated[dist.City.ID] = true
		}

		r.SubDistrict[i] = v

		if _, exist := subDistrictExist[v]; exist {
			o.Failure("sub_district_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("sub district"))
		} else {
			subDistrictExist[v] = 1
		}
	}

	r.CityStr = strings.Join(cityList, ",")

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"name.required":             util.ErrorInputRequired("name"),
		"business_type_id.required": util.ErrorInputRequired("business type"),
		"sls_man_id.required":       util.ErrorInputRequired("sales manager"),
		"area_id.required":          util.ErrorInputRequired("area"),
		"sub_district_id.required":  util.ErrorInputRequired("sub district"),
	}

	return messages
}
