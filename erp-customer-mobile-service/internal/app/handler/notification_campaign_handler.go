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

type NotificationCampaignHandler struct {
	Option                      global.HandlerOptions
	ServiceNotificationCampaign service.INotificationCampaignService
}

// URLMapping declare endpoint with handler function.
func (h *NotificationCampaignHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceNotificationCampaign = service.NewNotificationCampaignService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/list", h.GetHistoryCampaign, cMiddleware.Authorized("private"))
	r.POST("/update-read", h.UpdateRead, cMiddleware.Authorized("private"))
	r.POST("/unread", h.CountUnread, cMiddleware.Authorized("private"))
}

func (h NotificationCampaignHandler) GetHistoryCampaign(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationCampaignRequestGet

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceNotificationCampaign.GetHistoryCampaign(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h NotificationCampaignHandler) UpdateRead(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationCampaignRequestUpdateRead

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	e = h.ServiceNotificationCampaign.UpdateRead(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	return ctx.Serve(e)
}

func (h NotificationCampaignHandler) CountUnread(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req dto.NotificationCampaignRequestCountUnread

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	if req.Session, e = service.CustomerSession(ctx); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = h.ServiceNotificationCampaign.CountUnread(ctx.Request().Context(), &req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
