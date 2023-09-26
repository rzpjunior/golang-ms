// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package banner

import (
	"strconv"

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
	r.GET("", h.read, auth.Authorized("bnr_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("bnr_rdd"))
	r.POST("", h.create, auth.Authorized("bnr_crt"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("bnr_arc"))
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var (
		data                   []*model.Banner
		total, area, archetype int64
		areaStr, archetypeStr  string
	)

	if ctx.QueryParam("area") != "" {
		area, _ = common.Decrypt(ctx.QueryParam("area"))
		areaStr = strconv.Itoa(int(area))
	}

	if ctx.QueryParam("archetype") != "" {
		archetype, _ = common.Decrypt(ctx.QueryParam("archetype"))
		archetypeStr = strconv.Itoa(int(archetype))
	}

	if data, total, err = repository.GetBanners(ctx.RequestQuery(), areaStr, archetypeStr); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetBanner("id", id); e != nil {
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
