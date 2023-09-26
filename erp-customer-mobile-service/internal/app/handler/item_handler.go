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

type ItemHandler struct {
	Option      global.HandlerOptions
	ServiceItem service.IItemService
}

// URLMapping declare endpoint with handler function.
func (h *ItemHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceItem = service.NewItemService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.readPrivate, cMiddleware.Authorized("private"))
	r.POST("/detail", h.readPrivateDetail, cMiddleware.Authorized("private"))
	r.POST("/public", h.readPublic, cMiddleware.Authorized("public"))
	r.POST("/public/detail", h.readPublicDetail, cMiddleware.Authorized("public"))
	r.POST("/by-id", h.GetPrivateItemByListID, cMiddleware.Authorized("private"))
	r.POST("/last_finished_transaction", h.readLastFinTrans, cMiddleware.Authorized("private"))
}

func (h ItemHandler) readPrivate(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPrivateItemList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	// if req.Session, e = service.CustomerSession(ctx); e != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
	// 	return ctx.Serve(e)
	// }
	ctx.ResponseData, e = h.ServiceItem.GetPrivate(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemHandler) readPrivateDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.ItemDetailPrivateRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	// if req.Session, e = service.CustomerSession(ctx); e != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
	// 	return ctx.Serve(e)
	// }
	ctx.ResponseData, e = h.ServiceItem.GetPrivateDetail(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemHandler) readPublic(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetItemList

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	// if req.Session, e = service.CustomerSession(ctx); e != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
	// 	return ctx.Serve(e)
	// }
	ctx.ResponseData, e = h.ServiceItem.GetPublic(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemHandler) readPublicDetail(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.ItemDetailRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	// if req.Session, e = service.CustomerSession(ctx); e != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
	// 	return ctx.Serve(e)
	// }
	ctx.ResponseData, e = h.ServiceItem.GetPublicDetail(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h ItemHandler) GetPrivateItemByListID(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetPrivateItemByListID

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = h.ServiceItem.GetPrivateItemByListID(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

// get last finished transaction item list
func (h ItemHandler) readLastFinTrans(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RequestGetFinishedItems

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceItem.GetLastFinTrans(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}
