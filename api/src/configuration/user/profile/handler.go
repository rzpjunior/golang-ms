// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package profile

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	// r.GET("", h.get, auth.Authorized("usr_prf_rdl")) //TODO: Remove if unused
	// r.GET("/filter", h.getFilter, auth.Authorized("filter_rdl")) //TODO: Remove if unused
	r.GET("/detail", h.detail, auth.Authorized("usr_prf_rdd"))
	r.PUT("/detail", h.update, auth.Authorized("usr_prf_upd"))
	r.PUT("/password", h.updatePassword, auth.Authorized("usr_prf_upd_pas"))
}

/*
// get : function to get list of data TODO: Remove if unused
func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetStaffs(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}
 */

/*
func (h *Handler) getFilter(c echo.Context) (e error) { //TODO: Remove if unused
	ctx := c.(*cuxs.Context)
	data, total, e := repository.GetFilterStaff(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}
	return ctx.Serve(e)
}
*/

// detail : function to get detail of data based on id parameter
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var s *auth.SessionData
	if s, e = auth.UserSession(ctx); e == nil {
		if ctx.ResponseData, e = repository.GetStaff("id", s.Staff.ID); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update : function to update staff data based on user id
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Update(r)
		}
	}
	return ctx.Serve(e)
}

// updatePassword : function to update password of user
func (h *Handler) updatePassword(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updatePasswordRequest
	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = UpdatePassword(r)
		}
	}
	return ctx.Serve(e)
}
