// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package banner

import (
	"sort"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code             string    `json:"-"`
	Name             string    `json:"name" valid:"required"`
	AreasID          []string  `json:"areas" valid:"required"`
	ArchetypesID     []string  `json:"archetypes" valid:"required"`
	StartDate        string    `json:"start_date" valid:"required"`
	EndDate          string    `json:"end_date" valid:"required"`
	ImageUrl         string    `json:"image_url" valid:"required"`
	NavigationType   int8      `json:"navigation_type" valid:"required"`
	NavigationUrl    string    `json:"navigation_url"`
	TagProductID     string    `json:"tag_product_id"`
	ProductID        string    `json:"product_id"`
	Queue            int8      `json:"queue" valid:"required"`
	Note             string    `json:"note"`
	AreaStr          string    `json:"-"`
	ArchetypeStr     string    `json:"-"`
	StartTimestamp   time.Time `json:"-"`
	EndTimestamp     time.Time `json:"-"`
	ProductSectionID string    `json:"product_section_id"`

	TagProduct     *model.TagProduct     `json:"-"`
	Product        *model.Product        `json:"-"`
	ProductSection *model.ProductSection `json:"-"`

	Session *auth.SessionData
}

// Validate : function to validate request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                       error
		areaIDArr, archetypeIDArr []int
	)

	layout := "2006-01-02 15:04:05"
	prefix, err := util.CheckTable("banner")
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)

	r.Code = prefix + time.Now().Format("0106")

	if len(r.AreasID) > 0 {
		for i, v := range r.AreasID {
			var areaID int64

			if areaID, err = common.Decrypt(v); err != nil {
				o.Failure("areas_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("area"))
				continue
			}

			areaIDArr = append(areaIDArr, int(areaID))
		}
		sort.Ints(areaIDArr)

		for _, v := range areaIDArr {
			r.AreaStr += strconv.Itoa(v) + ","
		}
		r.AreaStr = r.AreaStr[:len(r.AreaStr)-1]
	}

	if len(r.ArchetypesID) > 0 {
		for i, v := range r.ArchetypesID {
			var archetypeID int64

			if archetypeID, err = common.Decrypt(v); err != nil {
				o.Failure("archetypes_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("archetype"))
				continue
			}

			archetypeIDArr = append(archetypeIDArr, int(archetypeID))
		}
		sort.Ints(archetypeIDArr)

		for _, v := range archetypeIDArr {
			r.ArchetypeStr += strconv.Itoa(v) + ","
		}
		r.ArchetypeStr = r.ArchetypeStr[:len(r.ArchetypeStr)-1]
	}

	if r.StartTimestamp, err = time.ParseInLocation(layout, r.StartDate, loc); err != nil {
		o.Failure("start_date.invalid", util.ErrorInvalidData("start date"))
	}

	if r.EndTimestamp, err = time.ParseInLocation(layout, r.EndDate, loc); err != nil {
		o.Failure("end_date.invalid", util.ErrorInvalidData("end date"))
	}

	if !r.StartTimestamp.IsZero() && r.StartTimestamp.Before(currentTime) {
		o.Failure("start_date.invalid", util.ErrorLater("start date", "current date"))
	}

	if !r.StartTimestamp.IsZero() && !r.EndTimestamp.IsZero() && (r.StartTimestamp.After(r.EndTimestamp) || r.StartTimestamp.Equal(r.EndTimestamp)) {
		o.Failure("start_date.invalid", util.ErrorLater("end date", "Start Date"))
	}

	if r.NavigationType == 1 {
		if r.NavigationUrl == "" {
			o.Failure("navigation_url.invalid", util.ErrorInputRequired("navigation url"))
			return o
		}
	} else if r.NavigationType == 2 {
		var tagProductID int64

		if r.TagProductID == "" {
			o.Failure("tag_product_id.invalid", util.ErrorSelectRequired("tag product"))
			return o
		}

		if tagProductID, err = common.Decrypt(r.TagProductID); err != nil {
			o.Failure("tag_product_id.invalid", util.ErrorInvalidData("tag product"))
			return o
		}

		if r.TagProduct, err = repository.ValidTagProduct(tagProductID); err != nil {
			o.Failure("tag_product_id.invalid", util.ErrorInvalidData("tag product"))
			return o
		}
	} else if r.NavigationType == 3 {
		var productID int64

		if r.ProductID == "" {
			o.Failure("product_id.invalid", util.ErrorSelectRequired("product"))
			return o
		}

		if productID, err = common.Decrypt(r.ProductID); err != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
			return o
		}

		if r.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
			return o
		}
	} else if r.NavigationType == 6 {
		var productSectionID int64
		if r.ProductSectionID == "" {
			o.Failure("product_section_id.invalid", util.ErrorSelectRequired("product section"))
			return o
		}

		if productSectionID, err = common.Decrypt(r.ProductSectionID); err != nil {
			o.Failure("product_section_id.invalid", util.ErrorInvalidData("product section"))
			return o
		}

		if r.ProductSection, err = repository.ValidProductSection(productSectionID); err != nil {
			o.Failure("product_section_id.invalid", util.ErrorInvalidData("product section"))
			return o
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":            util.ErrorInputRequired("name"),
		"areas.required":           util.ErrorSelectRequired("areas"),
		"archetypes.required":      util.ErrorSelectRequired("archetypes"),
		"start_date.required":      util.ErrorInputRequired("start date"),
		"end_date.required":        util.ErrorInputRequired("end date"),
		"image_url.required":       util.ErrorInputRequired("image"),
		"navigation_type.required": util.ErrorSelectRequired("Redirect to"),
		"queue.required":           util.ErrorInputRequired("queue"),
	}
}
