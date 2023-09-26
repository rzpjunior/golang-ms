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

type EdenPointHandler struct {
	Option           global.HandlerOptions
	ServiceEdenPoint service.IEdenPointService
}

// URLMapping declare endpoint with handler function.
func (h *EdenPointHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceEdenPoint = service.NewEdenPointService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/point_used", h.IsPointUsed, cMiddleware.Authorized("private"))
	r.POST("/potential", h.GetPotentialEdenPoint, cMiddleware.Authorized("private"))
	r.POST("/point_history", h.getPointHistoryList, cMiddleware.Authorized("private"))
	r.POST("/expiration", h.GetCustomerPointExpiration, cMiddleware.Authorized("private"))
}

func (h EdenPointHandler) IsPointUsed(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.IsPointUsedRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	res := &dto.IsPointUsedResponse{
		IsPointUsed: false,
	}
	res.IsPointUsed = h.ServiceEdenPoint.IsUsedEdenPoint(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData = res

	return ctx.Serve(e)
}

func (h EdenPointHandler) GetPotentialEdenPoint(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.GetPotentialEdenPointRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceEdenPoint.GetPotentialEdenPoint(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

func (h EdenPointHandler) getPointHistoryList(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.RequestGetPointHistory

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceEdenPoint.GetPointHistoryMobile(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h EdenPointHandler) GetCustomerPointExpiration(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.GetCustomerPointExpirationRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceEdenPoint.GetCustomerPointExpiration(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
