// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.signin)
	r.POST("/mobile", h.signinmobile)
	r.POST("/mobile/field_purchaser", h.signInMobileFieldPurchaser)
	r.GET("/me", h.me, cuxs.Authorized())
	r.GET("/checkmaintenance", h.check)
}

// signin endpoint to handle post http method.
func (h *Handler) signin(c echo.Context) (e error) {
	var r SignInRequest
	var sd *SessionData

	ctx := c.(*cuxs.Context)

	options := limiter.Options{
		TrustForwardHeader: true,
	}

	ipRequester := limiter.GetIP(ctx.Request(), options).String()
	r.IPAddress = ipRequester

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	if sd, e = Login(r.User); e != nil {
		return ctx.Serve(e)
	}
	ctx.Data(sd)
	return ctx.Serve(e)
}

// me endpoint untuk get sesion data yang lagi login.
func (h *Handler) me(c echo.Context) (e error) {
	var sd *SessionData
	ctx := c.(*cuxs.Context)
	// get current user dan data application menu
	if sd, e = UserSession(ctx); e == nil {
		ctx.Data(sd)
	}
	return ctx.Serve(e)
}

// check endpoint untuk check maintenance dino - dashboard.
func (h *Handler) check(c echo.Context) (e error) {
	var maintenance int
	o := orm.NewOrm()
	o.Using("read_only")
	o.Raw("SELECT value from config_app where attribute = 'maintenance_dino'").QueryRow(&maintenance)
	if maintenance == 1 {
		return echo.NewHTTPError(503, "server maintenance")
	}
	return

}

// signin endpoint to handle post http method.
func (h *Handler) signinmobile(c echo.Context) (e error) {
	var r SignInPackingRequest
	var sd *SessionData

	ctx := c.(*cuxs.Context)
	if e = ctx.Bind(&r); e == nil {
		if sd, e = Login(r.User); e == nil {
			ctx.Data(sd)
		}
	}
	return ctx.Serve(e)
}

// signin endpoint to mobile Field Purchaser Apps
func (h *Handler) signInMobileFieldPurchaser(c echo.Context) (e error) {
	var r SignInFieldPurchaserRequest
	var sd *SessionData

	ctx := c.(*cuxs.Context)

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	if sd, e = LoginFieldPurchaser(r); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(sd)

	return ctx.Serve(e)
}
