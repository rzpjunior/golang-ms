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

type ItemSectionHandler struct {
	Option             global.HandlerOptions
	ServiceItemSection service.IItemSectionService
}

// URLMapping declare endpoint with handler function.
func (h *ItemSectionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceItemSection = service.NewItemSectionService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.readPrivate, cMiddleware.Authorized("private"))
	r.POST("/public", h.read, cMiddleware.Authorized("public"))

	r.POST("/public/detail", h.readDetail, cMiddleware.Authorized("public"))
	r.POST("/detail", h.readPrivateDetail, cMiddleware.Authorized("public"))
}

func (h ItemSectionHandler) readPrivate(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPrivateItemSection

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceItemSection.GetPrivate(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemSectionHandler) read(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetItemSection

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceItemSection.GetPublic(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemSectionHandler) readDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetDetailItemSection

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceItemSection.GetPublicDetail(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemSectionHandler) readPrivateDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPrivateItemSectionDetail

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceItemSection.GetPrivateDetail(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
