// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package membership

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
	r.GET("/level", h.readLevel, auth.Authorized("mbr_lvl_rdl"))
	r.GET("/level/:id", h.detailLevel, auth.Authorized("mbr_lvl_rdd"))
	r.GET("/advantage", h.readAdvantage, auth.Authorized("mbr_adv_rdl"))
	r.GET("/advantage/:id", h.detailAdvantage, auth.Authorized("mbr_adv_rdd"))
	r.GET("/checkpoint", h.readCheckpoint, auth.Authorized("mbr_chp_rdl"))
	r.GET("/checkpoint/:id", h.detailCheckpoint, auth.Authorized("mbr_chp_rdd"))
}

// readLevel : function to get requested data of membership level based on parameters
func (h *Handler) readLevel(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var (
		data  []*model.MembershipLevel
		total int64
	)

	if data, total, err = repository.GetMembershipLevels(rq); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detailLevel : function to get detailed data of membership level by id
func (h *Handler) detailLevel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMembershipLevel("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// readAdvantage : function to get requested data of membership advantage based on parameters
func (h *Handler) readAdvantage(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var (
		data               []*model.MembershipAdvantage
		total, levelID     int64
		levelEnc, levelDec string
	)

	levelEnc = ctx.QueryParam("level.e")
	if levelEnc == "" {
		levelDec = ctx.QueryParam("level")
		levelID, _ = strconv.ParseInt(levelDec, 10, 64)
	} else {
		levelID, _ = common.Decrypt(levelEnc)
	}

	if data, total, err = repository.GetMembershipAdvantages(rq, levelID); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detailAdvantage : function to get detailed data of membership advantage by id
func (h *Handler) detailAdvantage(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMembershipAdvantage("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}

// readCheckpoint : function to get requested data of membership checkpoint based on parameters
func (h *Handler) readCheckpoint(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var (
		data  []*model.MembershipCheckpoint
		total int64
	)

	if data, total, err = repository.GetMembershipCheckpoints(rq); err != nil {
		err = echo.ErrNotFound
		return ctx.Serve(err)
	}

	ctx.Data(data, total)

	return ctx.Serve(err)
}

// detailCheckpoint : function to get detailed data of membership checkpoint by id
func (h *Handler) detailCheckpoint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if ctx.ResponseData, e = repository.GetMembershipCheckpoint("id", id); e != nil {
		e = echo.ErrNotFound
	}

	return ctx.Serve(e)
}
