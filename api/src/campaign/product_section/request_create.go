// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_section

import (
	"net/url"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code            string    `json:"-"`
	Name            string    `json:"name" valid:"required"`
	AreasID         []string  `json:"areas" valid:"required"`
	ArchetypesID    []string  `json:"archetypes" valid:"required"`
	StartAt         time.Time `json:"start_at" valid:"required"`
	EndAt           time.Time `json:"end_at" valid:"required"`
	BackgroundImage string    `json:"background_image"`
	ProductsID      []string  `json:"products" valid:"required"`
	Sequence        int8      `json:"sequence" valid:"required"`
	AreaStr         string    `json:"-"`
	ArchetypeStr    string    `json:"-"`
	ProductStr      string    `json:"-"`
	Status          int8      `json:"-"`
	Type            int8      `json:"type"`

	Session *auth.SessionData
}

// Validate : function to validate request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err     error
		isExist bool
	)

	// check code
	if r.Code, err = util.CheckTable("product_section"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	// check,decrypt array area_id and convert to string array
	for i, v := range r.AreasID {
		var areaID int64
		if areaID, err = common.Decrypt(v); err != nil {
			o.Failure("areas_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("area"))
			return o
		}
		areaIDStr := strconv.Itoa(int(areaID))
		if i != 0 {
			r.AreaStr += ","
		}
		r.AreaStr += areaIDStr
	}

	// check,decrypt array archetype_id and convert to string array
	for i, v := range r.ArchetypesID {
		var archetypeID int64
		if archetypeID, err = common.Decrypt(v); err != nil {
			o.Failure("archetypes_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("archetype"))
			return o
		}
		archetypeIDStr := strconv.Itoa(int(archetypeID))
		if i != 0 {
			r.ArchetypeStr += ","
		}
		r.ArchetypeStr += archetypeIDStr
	}

	// check sequence limitation
	if r.Sequence < 1 || r.Sequence > 5 {
		o.Failure("sequence.invalid", util.ErrorInvalidData("sequence"))
		return o
	}

	// only check if is prod recommendation is false
	if r.Type != 2 {
		if r.BackgroundImage == "" {
			o.Failure("background_image.required", util.ErrorInputRequired("background image"))
		}

		// check background image is url
		_, err = url.ParseRequestURI(r.BackgroundImage)
		if err != nil {
			o.Failure("background_image.invalid", util.ErrorInvalidData("background image"))
			return o
		}

		r.Type = 1
	} else {
		r.BackgroundImage = ""
	}

	// check,decrypt array product_id and convert to string array
	for i, v := range r.ProductsID {
		var productID int64
		if productID, err = common.Decrypt(v); err != nil {
			o.Failure("products_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			return o
		}

		// validate product
		_, err = repository.ValidProduct(productID)
		if err != nil {
			o.Failure("products_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			return o
		}

		productIDStr := strconv.Itoa(int(productID))
		if i != 0 {
			r.ProductStr += ","
		}
		r.ProductStr += productIDStr
	}

	// get glossary status
	statusGlossary, err := repository.GetGlossaryMultipleValue("table", "product_section", "attribute", "status", "value_name", "Draft")
	if err != nil {
		o.Failure("status.invalid", util.ErrorInvalidData("status"))
	}
	r.Status = statusGlossary.ValueInt

	// check start_at not greater than time.now
	if r.StartAt.Before(time.Now()) {
		o.Failure("start_at.invalid", util.ErrorGreater("start at", "time now"))
		return o
	}

	// check end_at not later than time.now and start_at
	if r.EndAt.Before(time.Now()) || r.EndAt.Before(r.StartAt) || r.EndAt.Equal(r.StartAt) {
		o.Failure("end_at.invalid", util.ErrorGreater("end at", "time now & start at"))
		return o
	}

	// only check if it is product recommendation
	if r.Type == 2 {
		// check if there are already data in between the date range
		if isExist, err = repository.CheckIsIntersect(r.Type, r.StartAt.Format("2006-01-02 15:04:05"), r.EndAt.Format("2006-01-02 15:04:05")); err != nil || isExist {
			o.Failure("start_at.invalid", util.ErrorIntersect("active", r.StartAt.Format("2006-01-02 15:04:05"), r.EndAt.Format("2006-01-02 15:04:05")))
			return o
		}
	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":       util.ErrorInputRequired("name"),
		"areas.required":      util.ErrorSelectRequired("areas"),
		"archetypes.required": util.ErrorSelectRequired("archetypes"),
		"start_at.required":   util.ErrorInputRequired("start at"),
		"end_at.required":     util.ErrorInputRequired("end at"),
		"sequence.required":   util.ErrorInputRequired("sequence"),
		"products.required":   util.ErrorInputRequired("products"),
	}
}
