package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type BannerHandler struct {
	Option        global.HandlerOptions
	ServiceBanner service.IBannerService
}

// URLMapping declare endpoint with handler function.
func (h *BannerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceBanner = service.NewBannerService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.readPrivate, cMiddleware.Authorized("private"))
	r.POST("/public", h.read, cMiddleware.Authorized("public"))
}

func (h BannerHandler) readPrivate(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPrivateBanner

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceBanner.GetPrivate(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h BannerHandler) read(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetBanner

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceBanner.GetPublic(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
