package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PushNotificationHandler struct {
	Option                   global.HandlerOptions
	ServicesPushNotification service.IPushNotificationService
}

// URLMapping implements router.RouteHandlers
func (h *PushNotificationHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPushNotification = service.NewPushNotificationService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("pnt_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("pnt_rdd"))
	r.POST("", h.Create, cMiddleware.Authorized("pnt_crt"))
	r.PUT("/:id", h.Update, cMiddleware.Authorized("pnt_upd"))
	r.PUT("/cancel/:id", h.Cancel, cMiddleware.Authorized("pnt_can"))
}

func (h PushNotificationHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	regionID := ctx.GetParamString("region_id")
	scheduledAtFrom := ctx.GetParamDate("scheduled_at_from")
	scheduledAtTo := ctx.GetParamDate("scheduled_at_to")

	var PushNotifications []*dto.PushNotificationResponse
	var total int64
	PushNotifications, total, err = h.ServicesPushNotification.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, regionID, scheduledAtFrom, scheduledAtTo)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(PushNotifications, total, page)

	return ctx.Serve(err)
}

func (h PushNotificationHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var PushNotification dto.PushNotificationResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	PushNotification, err = h.ServicesPushNotification.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = PushNotification

	return ctx.Serve(err)
}

func (h PushNotificationHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PushNotificationRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPushNotification.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PushNotificationHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.PushNotificationRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPushNotification.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PushNotificationHandler) Cancel(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.PushNotificationRequestCancel
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesPushNotification.Cancel(ctx.Request().Context(), req, id)

	return ctx.Serve(err)
}
