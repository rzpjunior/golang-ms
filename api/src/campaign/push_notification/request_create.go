// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package push_notification

import (
	"net/url"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold push notification request data
type createRequest struct {
	ID            string    `json:"-"`
	Code          string    `json:"-"`
	CampaignName  string    `json:"campaign_name" valid:"required"`
	Area          []int64   `json:"area" valid:"required"`
	Archetype     []int64   `json:"archetype" valid:"required"`
	RedirectTo    int8      `json:"redirect_to" valid:"required"`
	RedirectValue string    `json:"redirect_value"`
	Title         string    `json:"title" valid:"required"`
	Message       string    `json:"message" valid:"required"`
	PushNow       int8      `json:"push_now" valid:"required"`
	ScheduledAt   time.Time `json:"scheduled_at,omitempty"`

	Areas      []*model.Area      `json:"-"`
	Archetypes []*model.Archetype `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate push notification request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	// check code
	if c.Code, e = util.CheckTable("notification_campaign"); e != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	// checking areas array of id
	for _, a := range c.Area {
		idArea, e := common.Decrypt(a)
		if e != nil {
			o.Failure("area.invalid", util.ErrorInvalidData("area"))
			return o
		}
		area := &model.Area{
			ID: idArea,
		}
		e = area.Read("ID")
		if e != nil {
			o.Failure("area.invalid", util.ErrorInvalidData("area"))
			return o
		}

		c.Areas = append(c.Areas, area)
	}

	// checking archetypes array of id
	for _, a := range c.Archetype {
		idArchetypes, e := common.Decrypt(a)
		if e != nil {
			o.Failure("archetypes.invalid", util.ErrorInvalidData("archetypes"))
			return o
		}
		archetype := &model.Archetype{
			ID: idArchetypes,
		}
		e = archetype.Read("ID")
		if e != nil {
			o.Failure("archetypes_id.invalid", util.ErrorInvalidData("archetypes"))
			return o
		}

		c.Archetypes = append(c.Archetypes, archetype)

	}

	// checking redirect_to to glossary
	var redirectToGlossary *model.Glossary
	if redirectToGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "redirect_to", "value_int", c.RedirectTo); e != nil {
		o.Failure("redirect_to.invalid", util.ErrorInvalidData("redirect to"))
		return o
	}

	// check value name by redirect_to
	switch redirectToGlossary.ValueName {
	case "Product":
		id, e := common.Decrypt(c.RedirectValue)
		if e != nil {
			o.Failure("product_id.invalid", util.ErrorInvalidData("product"))
			return o
		}
		if _, e = repository.GetProduct("ID", id); e != nil {
			o.Failure("product.invalid", util.ErrorInvalidData("product"))
			return o
		}

		redirectValue := strconv.Itoa(int(id))
		c.RedirectValue = redirectValue
	case "Product Tag":
		id, e := common.Decrypt(c.RedirectValue)
		if e != nil {
			o.Failure("tag_product_id.invalid", util.ErrorInvalidData("tag product"))
			return o
		}
		if _, e = repository.GetProductTag("ID", id); e != nil {
			o.Failure("tag_product.invalid", util.ErrorInvalidData("tag product"))
			return o
		}

		redirectValue := strconv.Itoa(int(id))
		c.RedirectValue = redirectValue
	case "URL":
		_, e = url.ParseRequestURI(c.RedirectValue)
		if e != nil {
			o.Failure("url.invalid", util.ErrorInvalidData("url"))
			return o
		}
	}

	// checking push_now
	var pushNowGlossary *model.Glossary
	if pushNowGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", c.PushNow); e != nil {
		o.Failure("push_now.invalid", util.ErrorInvalidData("push now"))
		return o
	}

	// checking push now
	if pushNowGlossary.ValueName == "yes" {
		c.ScheduledAt = time.Now()
	} else {
		// check schedule is null
		if c.ScheduledAt.Year() == 1 {
			o.Failure("scheduled_at.invalid", util.ErrorInputRequired("scheduled at"))
			return o
		}
		// check greater than time now
		if c.ScheduledAt.Before(time.Now()) {
			o.Failure("scheduled_at.invalid", util.ErrorGreater("scheduled at", "time now"))
			return o
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"campaign_name.required": util.ErrorInputRequired("campaign name"),
		"area.required":          util.ErrorInputRequired("area"),
		"archetype.required":     util.ErrorInputRequired("archetype"),
		"redirect_to.required":   util.ErrorInputRequired("redirect to"),
		"title.required":         util.ErrorInputRequired("title"),
		"message.required":       util.ErrorInputRequired("message"),
		"push_now.required":      util.ErrorInputRequired("push now"),
	}
}
