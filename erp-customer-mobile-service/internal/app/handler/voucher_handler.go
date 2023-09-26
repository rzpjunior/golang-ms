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

type VoucherHandler struct {
	Option         global.HandlerOptions
	ServiceVoucher service.IVoucherService
}

// URLMapping declare endpoint with handler function.
func (h *VoucherHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceVoucher = service.NewVoucherService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.read, cMiddleware.Authorized("private"))
	r.POST("/detail", h.readDetail, cMiddleware.Authorized("private"))
	r.POST("/valid-promo", h.applyVoucher, cMiddleware.Authorized("private"))
	r.POST("/detail/items", h.GetVoucherItem, cMiddleware.Authorized("private"))
}

func (h VoucherHandler) read(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.VoucherRequestGet

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceVoucher.Get(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h VoucherHandler) readDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.VoucherRequestGetDetail

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceVoucher.GetDetail(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h VoucherHandler) applyVoucher(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.VoucherRequestApply

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceVoucher.ApplyVoucher(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h VoucherHandler) GetVoucherItem(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.VoucherRequestGetItemList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceVoucher.GetVoucherItem(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
