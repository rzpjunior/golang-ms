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

type SalesOrderHandler struct {
	Option            global.HandlerOptions
	ServiceSalesOrder service.ISalesOrderService
}

// URLMapping declare endpoint with handler function.
func (h *SalesOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceSalesOrder = service.NewSalesOrderService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.create, cMiddleware.Authorized("private"))
	r.POST("/update-cod", h.updateCOD, cMiddleware.Authorized("private"))
	r.POST("/feedback", h.getSalesOrderFeedback, cMiddleware.Authorized("private"))
	r.POST("/feedback/create", h.createSalesOrderFeedback, cMiddleware.Authorized("private"))
}

func (h SalesOrderHandler) create(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.CreateRequestSalesOrder

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceSalesOrder.Create(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h SalesOrderHandler) updateCOD(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.UpdateCodRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceSalesOrder.UpdateCOD(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h SalesOrderHandler) getSalesOrderFeedback(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.GetFeedback

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceSalesOrder.GetSalesOrderFeedback(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h SalesOrderHandler) createSalesOrderFeedback(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.CreateSalesFeedback

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceSalesOrder.CreateSalesOrderFeedback(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
