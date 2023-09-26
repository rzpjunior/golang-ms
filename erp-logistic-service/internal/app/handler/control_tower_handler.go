package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ControlTowerHandler struct {
	Option               global.HandlerOptions
	ServicesControlTower service.IControlTowerService
}

// URLMapping implements router.RouteHandlers
func (h *ControlTowerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesControlTower = service.NewServiceControlTower()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.GetDRS, cMiddleware.Authorized("ctrl_twr_rdl"))
	r.POST("", h.GetCourier, cMiddleware.Authorized("ctrl_twr_rdl"))
	r.GET("/:id", h.DetailDRS, cMiddleware.Authorized("ctrl_twr_rdd"))
	r.POST("/:id", h.DetailCourier, cMiddleware.Authorized("ctrl_twr_rdd"))
	r.PUT("/cancel/:id", h.CancelDRS, cMiddleware.Authorized("ctrl_twr_can"))
	r.PUT("/cancel/item/:id", h.CancelItem, cMiddleware.Authorized("ctrl_twr_can"))
}

func (h ControlTowerHandler) GetDRS(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	// mandatory fields, return nothing if not exist
	siteID := ctx.GetParamString("site_id")
	if siteID == "" {
		return ctx.Serve(err)
	}

	startDeliveryDate := ctx.GetParamDate("start_delivery_date")
	if !timex.IsValid(startDeliveryDate) {
		return ctx.Serve(err)
	}

	endDeliveryDate := ctx.GetParamDate("end_delivery_date")
	if !timex.IsValid(endDeliveryDate) {
		return ctx.Serve(err)
	}

	// optional fields
	vendorID := ctx.GetParamString("vendor_id")
	courierID := ctx.GetParamString("courier_id")
	statusIDs := ctx.GetParamArrayInt("status_id_in")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")

	var items []dto.ControlTowerGetDRSResponse
	var total int64
	items, total, err = h.ServicesControlTower.GetDRS(ctx.Request().Context(), dto.ControlTowerGetDRSRequest{
		Offset:            page.Start,
		Limit:             page.Limit,
		OrderBy:           orderBy,
		SiteID:            siteID,
		StartDeliveryDate: startDeliveryDate,
		EndDeliveryDate:   endDeliveryDate,
		CourierVendorID:   vendorID,
		CourierID:         courierID,
		StatusIDs:         statusIDs,
		Search:            search,
	})

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h ControlTowerHandler) GetCourier(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ControlTowerGetCourierRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesControlTower.GetCourier(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h ControlTowerHandler) DetailDRS(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var deliveryRunSheet dto.ControlTowerGetDRSDetailResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	deliveryRunSheet, err = h.ServicesControlTower.GetDRSDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = deliveryRunSheet

	return ctx.Serve(err)
}

func (h ControlTowerHandler) DetailCourier(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesControlTower.GetCourierDetail(ctx.Request().Context(), id)

	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h ControlTowerHandler) CancelDRS(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ControlTowerCancelDRSRequest

	if req.DeliveryRunSheetID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesControlTower.CancelDRS(ctx.Request().Context(), req)

	return ctx.Serve(err)
}

func (h ControlTowerHandler) CancelItem(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ControlTowerCancelItemRequest

	if req.DeliveryRunSheetItemID, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesControlTower.CancelItem(ctx.Request().Context(), req)

	return ctx.Serve(err)
}
