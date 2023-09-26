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

type AddressHandler struct {
	Option         global.HandlerOptions
	ServiceAddress service.IAddressService
}

// URLMapping declare endpoint with handler function.
func (h *AddressHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceAddress = service.NewAddressService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/list", h.read, cMiddleware.Authorized("private"))
	r.POST("/create", h.create, cMiddleware.Authorized("private"))
	r.POST("/update", h.update, cMiddleware.Authorized("private"))
	r.POST("/default", h.setDefault, cMiddleware.Authorized("private"))
	r.POST("/delete", h.delete, cMiddleware.Authorized("private"))
}

func (h AddressHandler) read(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetAddressList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceAddress.Get(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h AddressHandler) create(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.CreateAddressRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	e = h.ServiceAddress.Create(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h AddressHandler) update(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.UpdateAddressRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	e = h.ServiceAddress.Update(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h AddressHandler) setDefault(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.SetDefaultAddressRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	e = h.ServiceAddress.SetDefault(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h AddressHandler) delete(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.DeleteAddressRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	e = h.ServiceAddress.Delete(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
