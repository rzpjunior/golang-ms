// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package edenpoint

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("ep_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("ep_rdd"))
	r.POST("", h.create, auth.Authorized("ep_crt"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("ep_arc"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var (
		areaID, archetypeID, customerTagID int64
		period                             []string
	)

	if ctx.QueryParam("area_id") != "" {
		areaID, err = common.Decrypt(ctx.QueryParam("area_id"))
	}

	if ctx.QueryParam("archetype_id") != "" {
		archetypeID, err = common.Decrypt(ctx.QueryParam("archetype_id"))
	}

	if ctx.QueryParam("customer_tag_id") != "" {
		customerTagID, err = common.Decrypt(ctx.QueryParam("customer_tag_id"))
	}

	if ctx.QueryParam("period") != "" {
		dateParam := ctx.QueryParam("period")
		period = strings.Split(dateParam, "|")
	}

	var data []*model.EdenPointCampaign
	var total int64

	if data, total, err = repository.GetEdenPointCampaigns(rq, areaID, archetypeID, customerTagID, period); err != nil {
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var campaignID int64

	if campaignID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetEdenPointCampaign("id", campaignID); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}

	return ctx.Serve(e)
}

// archive : function to change status of data into archive status
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Archive(r)

	return ctx.Serve(e)
}
