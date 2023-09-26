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

type NotificationTransactionHandler struct {
	Option                         global.HandlerOptions
	ServiceNotificationTransaction service.INotificationTransactionService
}

// URLMapping declare endpoint with handler function.
func (h *NotificationTransactionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceNotificationTransaction = service.NewNotificationTransactionService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/list", h.GetHistoryTransaction, cMiddleware.Authorized("private"))
	r.POST("/update-read", h.UpdateRead, cMiddleware.Authorized("private"))
	r.POST("/unread", h.CountUnread, cMiddleware.Authorized("private"))
}

func (h NotificationTransactionHandler) GetHistoryTransaction(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationTransactionRequestGet

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceNotificationTransaction.GetHistoryTransaction(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h NotificationTransactionHandler) UpdateRead(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationTransactionRequestUpdateRead

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	e = h.ServiceNotificationTransaction.UpdateRead(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h NotificationTransactionHandler) CountUnread(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationTransactionRequestCountUnread

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceNotificationTransaction.CountUnread(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
