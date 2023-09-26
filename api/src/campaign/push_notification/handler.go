// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package push_notification

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pnt_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pnt_rdd"))
	r.POST("", h.create, auth.Authorized("pnt_crt"))
	r.PUT("/:id", h.update, auth.Authorized("pnt_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("pnt_cnl"))
}

// read : function to get all data
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		data        []*model.NotificationCampaign
		total, area int64
		areaStr     string
	)

	if ctx.QueryParam("area") != "" {
		area, _ = common.Decrypt(ctx.QueryParam("area"))
		areaStr = strconv.Itoa(int(area))
	}

	if data, total, e = repository.GetNotificationCampaigns(ctx.RequestQuery(), areaStr); e != nil {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	data, e := repository.GetNotificationCampaign("id", id)
	if e != nil {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}
	ctx.ResponseData = data

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			var data *model.NotificationCampaign
			data, e = Save(r)
			if e != nil {
				return ctx.Serve(e)
			}
			ctx.ResponseData = data

			areaArr := strings.Split(data.Area, ",")
			archetypeArr := strings.Split(data.Archetype, ",")
			var areaIDs, archetypeIDs []int64

			for _, i := range areaArr {
				id, _ := strconv.Atoi(i)
				areaIDs = append(areaIDs, int64(id))
			}
			for _, i := range archetypeArr {
				id, _ := strconv.Atoi(i)
				archetypeIDs = append(archetypeIDs, int64(id))
			}

			// checking push_now
			var pushNowGlossary *model.Glossary
			if pushNowGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", data.PushNow); e != nil {
				return ctx.Serve(e)
			}
			// checking push now
			if pushNowGlossary.ValueName == "yes" {
				data.PushNowStatus = true
			} else {
				data.PushNowStatus = false
			}

			data, e = repository.GetNotificationCampaign("id", data.ID)
			if e != nil {
				e = echo.ErrNotFound
				return ctx.Serve(e)
			}

			if data.PushNowStatus {
				messageNotif := &util.MessageNotificationCampaign{
					ID:             common.Encrypt(data.ID),
					Code:           data.Code,
					CampaignName:   data.CampaignName,
					Area:           areaIDs,
					Archetype:      archetypeIDs,
					RedirectTo:     data.RedirectTo,
					RedirectToName: data.RedirectToName,
					RedirectValue:  data.RedirectValue,
					Title:          data.Title,
					Message:        data.Message,
					ServerKey:      util.CampaignServerKeyFireBase,
				}

				if e = util.PostModelNotificationCampaign(messageNotif); e != nil {
					return ctx.Serve(e)
				}
			}
		}
	}

	return ctx.Serve(e)
}

//update : function to update data based on input
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	var data *model.NotificationCampaign
	data, e = Update(r)
	if e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData = data

	areaArr := strings.Split(data.Area, ",")
	archetypeArr := strings.Split(data.Archetype, ",")
	var areaIDs, archetypeIDs []int64

	for _, i := range areaArr {
		id, _ := strconv.Atoi(i)
		areaIDs = append(areaIDs, int64(id))
	}
	for _, i := range archetypeArr {
		id, _ := strconv.Atoi(i)
		archetypeIDs = append(archetypeIDs, int64(id))
	}

	// checking push_now
	var pushNowGlossary *model.Glossary
	if pushNowGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", data.PushNow); e != nil {
		return ctx.Serve(e)
	}
	// checking push now
	if pushNowGlossary.ValueName == "yes" {
		data.PushNowStatus = true
	} else {
		data.PushNowStatus = false
	}

	data, e = repository.GetNotificationCampaign("id", data.ID)
	if e != nil {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}

	if data.PushNowStatus {
		messageNotif := &util.MessageNotificationCampaign{
			ID:             common.Encrypt(data.ID),
			Code:           data.Code,
			CampaignName:   data.CampaignName,
			Area:           areaIDs,
			Archetype:      archetypeIDs,
			RedirectTo:     data.RedirectTo,
			RedirectToName: data.RedirectToName,
			RedirectValue:  data.RedirectValue,
			Title:          data.Title,
			Message:        data.Message,
			ServerKey:      util.CampaignServerKeyFireBase,
		}

		if e = util.PostModelNotificationCampaign(messageNotif); e != nil {
			return ctx.Serve(e)
		}
	}

	ctx.ResponseData = data
	return ctx.Serve(e)
}

//cancel : function to cancel data based on input
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Cancel(r)
	return ctx.Serve(e)
}
