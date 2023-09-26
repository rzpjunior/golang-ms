// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_section

import (
	"strconv"
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
	r.GET("/product", h.readProduct, auth.Authorized("psc_rdl"))
	r.GET("", h.read, auth.Authorized("psc_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("psc_rdd"))
	r.POST("", h.create, auth.Authorized("psc_crt"))
	r.PUT("/:id", h.update, auth.Authorized("psc_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("psc_arc"))
}

// readProduct : function to get product requested data based on parameters
func (h *Handler) readProduct(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)

	var (
		data     []*model.ProductSectionItem
		total    int64
		category string
	)

	if ctx.QueryParam("category") != "" {
		category = ctx.QueryParam("category")
	}

	if data, total, err = repository.GetProductSectionItem(ctx.RequestQuery(), category); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)
	return ctx.Serve(err)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var (
		data                   []*model.ProductSection
		total, area, archetype int64
		areaStr, archetypeStr  string
		statusArr              []int
	)

	if ctx.QueryParam("area") != "" {
		listArea := strings.Split(ctx.QueryParam("area"), ",")
		for _, v := range listArea {
			area, _ = common.Decrypt(v)
			areaStr += strconv.Itoa(int(area)) + ","
		}
		areaStr = strings.TrimSuffix(areaStr, ",")
	}

	if ctx.QueryParam("archetype") != "" {
		listArchetype := strings.Split(ctx.QueryParam("archetype"), ",")
		for _, v := range listArchetype {
			archetype, _ = common.Decrypt(v)
			archetypeStr += strconv.Itoa(int(archetype)) + ","
		}
		archetypeStr = strings.TrimSuffix(archetypeStr, ",")
	}

	if ctx.QueryParam("status") != "" {
		listStatus := strings.Split(ctx.QueryParam("status"), ",")
		for _, v := range listStatus {
			statusInt, _ := strconv.Atoi(v)
			statusArr = append(statusArr, statusInt)
		}
	}

	if data, total, err = repository.GetProductSections(ctx.RequestQuery(), areaStr, archetypeStr, statusArr); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, err = ctx.Decrypt("id"); err != nil {
		return ctx.Serve(err)
	}

	if ctx.ResponseData, err = repository.GetProductSection("id", id); err != nil {
		err = echo.ErrNotFound
	}

	return ctx.Serve(err)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, err = auth.UserSession(ctx); err != nil {
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&r); err == nil {
		ctx.ResponseData, err = Save(r)
	}

	return ctx.Serve(err)
}

// update : function to update requested data based on parameters
func (h *Handler) update(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, err = auth.UserSession(ctx); err != nil {
		return ctx.Serve(err)
	}

	if r.ID, err = ctx.Decrypt("id"); err != nil {
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&r); err != nil {
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = Update(r)

	return ctx.Serve(err)
}

// archive : function to change status of data into archive status
func (h *Handler) archive(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, err = auth.UserSession(ctx); err != nil {
		return ctx.Serve(err)
	}

	if r.ID, err = ctx.Decrypt("id"); err != nil {
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&r); err != nil {
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = Archive(r)

	return ctx.Serve(err)
}
