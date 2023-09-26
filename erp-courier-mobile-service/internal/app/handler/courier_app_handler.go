package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CourierAppHandler struct {
	Option             global.HandlerOptions
	ServicesCourierApp service.ICourierAppService
}

// URLMapping implements router.RouteHandlers
func (h *CourierAppHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCourierApp = service.NewServiceCourierApp()

	cMiddleware := middleware.NewMiddleware()

	r.POST("/login", h.Login)
	r.POST("/courier/log", h.CourierLog, cMiddleware.AuthorizedCourierApp())
	r.PUT("/activate/emergency", h.ActivateEmergency, cMiddleware.AuthorizedCourierApp())
	r.PUT("/deactivate/emergency", h.DeactivateEmergency, cMiddleware.AuthorizedCourierApp())

	// Get Data
	r.GET("", h.Get, cMiddleware.AuthorizedCourierApp())
	r.GET("/:id", h.Detail, cMiddleware.AuthorizedCourierApp())
	r.POST("/scan", h.ScanDetail, cMiddleware.AuthorizedCourierApp())

	// Assignment Creation
	r.POST("/self-assign/scan", h.Scan, cMiddleware.AuthorizedCourierApp())
	r.POST("/self-assign/:id", h.SelfAssign, cMiddleware.AuthorizedCourierApp())

	// Actions
	//// Delivery
	r.PUT("/start/delivery/:id", h.StartDelivery, cMiddleware.AuthorizedCourierApp())
	r.PUT("/success/delivery/:id", h.SuccessDelivery, cMiddleware.AuthorizedCourierApp())
	r.PUT("/postpone/delivery/:id", h.PostponeDelivery, cMiddleware.AuthorizedCourierApp())
	r.PUT("/fail/delivery/:id", h.FailDelivery, cMiddleware.AuthorizedCourierApp())
	r.POST("/status/delivery/:id", h.StatusDelivery, cMiddleware.AuthorizedCourierApp())
	//// Return
	// r.POST("/return/:id", h.CreateReturn, cMiddleware.AuthorizedCourierApp())
	// r.PUT("/return/:id", h.EditReturn, cMiddleware.AuthorizedCourierApp())
	// r.DELETE("/return/:id", h.DeleteReturn, cMiddleware.AuthorizedCourierApp())

	// Upload Picture
	r.POST("/merchant/delivery/log/:id", h.CreateMerchantDeliveryLog, cMiddleware.AuthorizedCourierApp())

	// Get Glossary
	r.GET("/glossary", h.GetGlossary, cMiddleware.AuthorizedCourierApp())
}

func (h CourierAppHandler) Login(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.LoginRequest

	timezone := ctx.Request().Header.Get("Timezone")
	if timezone != "" {
		req.Timezone = timezone
	} else {
		req.Timezone = "Asia/Jakarta"
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.Login(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) CourierLog(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CreateCourierLogRequest

	req.CourierID = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesCourierApp.CreateCourierLog(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) ActivateEmergency(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppActivateEmergencyRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	ctx.ResponseData, err = h.ServicesCourierApp.ActivateEmergency(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) DeactivateEmergency(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.CourierAppDeactivateEmergencyRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	ctx.ResponseData, err = h.ServicesCourierApp.DeactivateEmergency(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppGetRequest

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	req.Offset = page.Start
	req.Limit = page.Limit
	req.OrderBy = ctx.GetParamString("order_by")

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	req.StartDeliveryDate = ctx.GetParamDate("start_delivery_date")
	req.EndDeliveryDate = ctx.GetParamDate("end_delivery_date")
	req.StepType = ctx.GetParamInt("step_type")
	req.StatusIDs = ctx.GetParamArrayInt("status_id_in")
	req.Search = ctx.GetParamString("search")
	req.SearchSalesOrderCode = ctx.GetParamString("search_sales_order_code")
	var items []dto.CourierAppGetResponse
	var total int64
	items, total, err = h.ServicesCourierApp.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h CourierAppHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var deliveryRunSheetItem dto.CourierAppDetailResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	deliveryRunSheetItem, err = h.ServicesCourierApp.Detail(ctx.Request().Context(), id, ctx.Request().Context().Value(constants.KeyCourierID).(string))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = deliveryRunSheetItem

	return ctx.Serve(err)
}

func (h CourierAppHandler) ScanDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppScanDetailRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.ScanDetail(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) Scan(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppScanRequest

	req.CourierSiteId = ctx.Request().Context().Value(constants.KeySiteID).(string)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.Scan(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) SelfAssign(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppSelfAssignRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	req.CourierSiteId = ctx.Request().Context().Value(constants.KeySiteID).(string)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.SelfAssign(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) StartDelivery(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppStartDeliveryRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	req.CourierSiteId = ctx.Request().Context().Value(constants.KeySiteID).(string)

	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.StartDelivery(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) SuccessDelivery(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppSuccessDeliveryRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	req.CourierSiteId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.SuccessDelivery(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) PostponeDelivery(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppPostponeDeliveryRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)
	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.PostponeDelivery(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) FailDelivery(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppFailDeliveryRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.FailDelivery(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) StatusDelivery(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppStatusDeliveryRequest

	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.StatusDelivery(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) CreateMerchantDeliveryLog(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppCreateMerchantDeliveryLogRequest

	req.CourierId = ctx.Request().Context().Value(constants.KeyCourierID).(string)

	if req.Id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesCourierApp.CreateMerchantDeliveryLog(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h CourierAppHandler) GetGlossary(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.CourierAppGetGlossaryRequest

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req.Table = ctx.GetParamString("table")
	req.Attribute = ctx.GetParamString("attribute")
	req.ValueInt = ctx.GetParamInt("value_int")
	req.ValueName = ctx.GetParamString("value_name")

	glossaries, total, err := h.ServicesCourierApp.GetGlossary(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(glossaries, total, page)

	return ctx.Serve(err)
}
