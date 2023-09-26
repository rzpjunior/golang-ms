// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_group

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"strconv"
	"strings"
)

type updateRequest struct {
	ID             int64    `json:"-" valid:"required"`
	Name           string   `json:"name" valid:"required"`
	BusinessTypeID string   `json:"business_type_id" valid:"required"`
	SlsManID       string   `json:"sls_man_id" valid:"required"`
	SubDistrict    []string `json:"sub_district_id" valid:"required"`
	SubDistrictStr string   `json:"-"`
	CityStr        string   `json:"-"`

	SalesGroup     *model.SalesGroup     `json:"-"`
	BusinessType   *model.BusinessType   `json:"-"`
	SlsMan         *model.Staff          `json:"-"`
	SalesGroupItem *model.SalesGroupItem `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")
	var err error
	var businessTypeID, slsManID int64
	var subDistrictExist = make(map[string]int8)
	var duplicated = make(map[int64]bool)
	var cityList []string

	if c.SalesGroup, err = repository.ValidSalesGroup(c.ID); err != nil {
		o.Failure("sales_group.invalid", util.ErrorInvalidData("sales_group"))
	}

	if c.SalesGroup.Status != 1 {
		o.Failure("id.invalid", util.ErrorActive("status"))
	}

	filter := map[string]interface{}{"name": c.Name, "status": int8(1)}
	exclude := map[string]interface{}{"id": c.ID, "status": int8(2)}
	if _, countName, err := repository.CheckSalesGroupData(filter, exclude); err != nil {
		o.Failure("name.invalid", util.ErrorInvalidData("name"))
	} else if countName > 0 {
		o.Failure("name.unique", util.ErrorUnique("name"))
	}

	if businessTypeID, err = common.Decrypt(c.BusinessTypeID); err != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	if c.BusinessType, err = repository.ValidBusinessType(businessTypeID); err != nil {
		o.Failure("business_type_id.invalid", util.ErrorInvalidData("business type"))
	}

	if c.BusinessType.Status != 1 {
		o.Failure("business_type_id.inactive", util.ErrorActive("business type"))
	}

	if slsManID, err = common.Decrypt(c.SlsManID); err != nil {
		o.Failure("sls_man_id.invalid", util.ErrorInvalidData("sales manager"))
	}

	if c.SlsMan, err = repository.ValidStaff(slsManID); err != nil {
		o.Failure("sls_man_id.invalid", util.ErrorInvalidData("sales manager"))
	}

	for i, v := range c.SubDistrict {
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

		c.SubDistrict[i] = v

		if _, exist := subDistrictExist[v]; exist {
			o.Failure("sub_district_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("sub district"))
		} else {
			subDistrictExist[v] = 1
		}
	}

	c.CityStr = strings.Join(cityList, ",")

	return o
}

// Messages : function to return error messages after validation
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":             util.ErrorInputRequired("name"),
		"business_type_id.required": util.ErrorInputRequired("business type"),
		"sls_man_id.required":       util.ErrorInputRequired("sales manager"),
		"area_id.required":          util.ErrorInputRequired("area"),
		"sub_district_id.required":  util.ErrorInputRequired("sub district"),
	}
}
