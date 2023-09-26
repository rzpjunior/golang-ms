// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenpoint

import (
	"sort"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code               string    `json:"-"`
	Name               string    `json:"name" valid:"required"`
	AreasID            []string  `json:"areas"`
	ArchetypesID       []string  `json:"archetypes"`
	CustomerTagsID     []string  `json:"customer_tags"`
	CampaignFilterType int8      `json:"campaign_filter_type" valid:"required"`
	StartDate          string    `json:"start_date" valid:"required"`
	EndDate            string    `json:"end_date" valid:"required"`
	ImageUrl           string    `json:"image_url" valid:"required"`
	Multiplier         int8      `json:"multiple" valid:"required"`
	Note               string    `json:"note"`
	AreaStr            string    `json:"-"`
	ArchetypeStr       string    `json:"-"`
	CustomerTagStr     string    `json:"-"`
	StartTimestamp     time.Time `json:"-"`
	EndTimestamp       time.Time `json:"-"`

	Session *auth.SessionData
}

// Validate : function to validate request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                                         error
		areaIDArr, archetypeIDArr, customerTagIDArr []int
	)

	layout := "2006-01-02 15:04:05"
	prefix, err := util.CheckTable("eden_point_campaign")
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)

	r.Code = prefix + time.Now().Format("0106")

	if r.CampaignFilterType == 1 {
		if len(r.AreasID) == 0 {
			o.Failure("areas.required", util.ErrorInputRequired("area"))
		} else {
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

		if len(r.ArchetypesID) == 0 {
			o.Failure("archetypes.required", util.ErrorInputRequired("archetype"))
		} else {
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
	} else if r.CampaignFilterType == 2 {
		if len(r.CustomerTagsID) == 0 {
			o.Failure("customer_tags.required", util.ErrorInputRequired("customer tag"))
		} else {
			for i, v := range r.CustomerTagsID {
				var customerTagID int64

				if customerTagID, err = common.Decrypt(v); err != nil {
					o.Failure("customer_tags_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("customer tag"))
					continue
				}

				customerTagIDArr = append(customerTagIDArr, int(customerTagID))
			}
			sort.Ints(customerTagIDArr)

			for _, v := range customerTagIDArr {
				r.CustomerTagStr += strconv.Itoa(v) + ","
			}
			r.CustomerTagStr = r.CustomerTagStr[:len(r.CustomerTagStr)-1]
		}
	}

	if r.Multiplier < 2 {
		o.Failure("multiple.invalid", util.ErrorGreater("multiple", "1"))
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

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":                 util.ErrorInputRequired("name"),
		"campaign_filter_type.required": util.ErrorInputRequired("campaign filter type"),
		"start_date.required":           util.ErrorInputRequired("start date"),
		"end_date.required":             util.ErrorInputRequired("end date"),
		"image_url.required":            util.ErrorInputRequired("image"),
		"multiple.required":             util.ErrorInputRequired("multiple"),
	}
}
