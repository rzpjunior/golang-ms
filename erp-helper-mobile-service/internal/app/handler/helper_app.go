package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type HelperAppHandler struct {
	Option           global.HandlerOptions
	ServiceHelperApp service.IHelperAppService
}

// URLMapping implements router.RouteHandlers
func (h *HelperAppHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceHelperApp = service.NewServiceHelperApp()

	cMiddleware := middleware.NewMiddleware()

	// PICKER
	r.POST("/login", h.Login)
	r.GET("/picker", h.GetPickingOrder, cMiddleware.AuthorizedHelperMobile())
	r.GET("/picker/widget", h.PickerWidget, cMiddleware.AuthorizedHelperMobile())
	r.GET("/picker/detail", h.GetPickingOrderProducts, cMiddleware.AuthorizedHelperMobile())
	r.POST("/picker/detail/start", h.StartPickingOrder, cMiddleware.AuthorizedHelperMobile())
	r.GET("/picker/detail/sales-order", h.GetPickingOrderProductsSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.POST("/picker/detail/sales-order/submit", h.SubmitPicking, cMiddleware.AuthorizedHelperMobile())
	r.GET("/picker/detail/review", h.GetSalesOrderList, cMiddleware.AuthorizedHelperMobile())
	// USED FOR PICKER HISTORY AS WELL
	r.GET("/picker/sales-order/detail", h.GetSalesOrderDetail, cMiddleware.AuthorizedHelperMobile())
	r.POST("/picker/detail/review/submit", h.SubmitSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.GET("/picker/history", h.History, cMiddleware.AuthorizedHelperMobile())

	// SPV
	r.GET("/spv", h.SPVGetSalesOrderList, cMiddleware.AuthorizedHelperMobile())
	r.GET("/spv/widget", h.SPVWidget, cMiddleware.AuthorizedHelperMobile())
	r.GET("/spv/detail", h.SPVGetSalesOrderDetail, cMiddleware.AuthorizedHelperMobile())
	r.POST("/spv/detail/reject", h.SPVRejectSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.POST("/spv/detail/accept", h.SPVAcceptSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.POST("/spv/monitoring", h.SPVGetWrtMonitoring, cMiddleware.AuthorizedHelperMobile())
	r.POST("/spv/monitoring/detail", h.SPVGetWrtMonitoringDetail, cMiddleware.AuthorizedHelperMobile())

	// CHECKER
	r.GET("/checker", h.CheckerGetSalesOrderList, cMiddleware.AuthorizedHelperMobile())
	r.GET("/checker/widget", h.CheckerWidget, cMiddleware.AuthorizedHelperMobile())
	r.GET("/checker/detail", h.CheckerGetSalesOrderDetail, cMiddleware.AuthorizedHelperMobile())
	r.POST("/checker/detail/start", h.CheckerStartChecking, cMiddleware.AuthorizedHelperMobile())
	r.POST("/checker/detail/submit", h.CheckerSubmitChecking, cMiddleware.AuthorizedHelperMobile())
	r.POST("/checker/detail/reject", h.CheckerRejectSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.GET("/checker/detail/delivery-koli", h.CheckerGetDeliveryKoli, cMiddleware.AuthorizedHelperMobile())
	r.POST("/checker/detail/accept", h.CheckerAcceptSalesOrder, cMiddleware.AuthorizedHelperMobile())
	r.GET("/checker/history", h.CheckerHistory, cMiddleware.AuthorizedHelperMobile())
	// TODO TO BE REMOVED
	r.GET("/checker/history/detail", h.CheckerHistoryDetail, cMiddleware.AuthorizedHelperMobile())
}

func (h HelperAppHandler) Login(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppLoginRequest

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

	ctx.ResponseData, err = h.ServiceHelperApp.Login(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) GetPickingOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.HelperAppGetPickingOrderRequest{
		Limit:        page.PerPage,
		Offset:       page.Page - 1,
		LocationCode: ctx.GetParamString("locncode"),
		SopNumber:    ctx.GetParamString("sopnumbe"),
		DocNumber:    ctx.GetParamString("docnumbr"),
		ItemNumber:   ctx.GetParamString("itemnmbr"),
		HelperId:     ctx.Request().Context().Value(constants.KeyUserID).(string),
		Status:       int8(ctx.GetParamInt("status")),
		CustomerName: ctx.GetParamString("custname"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.GetPickingOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) PickerWidget(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppPickerWidgetRequest{
		HelperId: ctx.Request().Context().Value(constants.KeyUserID).(string),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.PickerWidget(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) GetPickingOrderProducts(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppGetPickingOrderProductsRequest{
		DocNumber:      ctx.GetParamString("doc_number"),
		ItemNameSearch: ctx.GetParamString("item_name_search"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.GetPickingOrderProducts(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) GetPickingOrderProductsSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	docNumber := ctx.GetParamString("doc_number")
	itemNumber := ctx.GetParamString("item_number")

	ctx.ResponseData, err = h.ServiceHelperApp.GetPickingOrderProductsSalesOrder(ctx.Request().Context(), docNumber, itemNumber)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) StartPickingOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppDocNumberRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceHelperApp.StartPickingOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SubmitPicking(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppSubmitPickingRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceHelperApp.SubmitPicking(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) GetSalesOrderList(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	docNumber := ctx.GetParamString("doc_number")

	ctx.ResponseData, err = h.ServiceHelperApp.GetSalesOrderPicking(ctx.Request().Context(), docNumber)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) GetSalesOrderDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	sopNumber := ctx.GetParamString("sop_number")

	ctx.ResponseData, err = h.ServiceHelperApp.GetSalesOrderPickingDetail(ctx.Request().Context(), sopNumber)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SubmitSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppSubmitSalesOrderRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.PickerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.SubmitSalesOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) History(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.HelperAppHistoryRequest{
		Limit:        page.PerPage,
		Offset:       page.Page - 1,
		SopNumber:    ctx.GetParamString("sop_number"),
		PickerId:     ctx.Request().Context().Value(constants.KeyUserID).(string),
		CustomerName: ctx.GetParamString("merchant_name"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.History(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVGetSalesOrderList(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var req dto.HelperAppGetSalesOrderToCheckRequest

	req = dto.HelperAppGetSalesOrderToCheckRequest{
		Offset:       page.Offset,
		Limit:        page.Limit,
		SiteId:       ctx.Request().Context().Value(constants.KeySiteID).(string),
		SopNumber:    ctx.GetParamString("sop_number"),
		CustomerName: ctx.GetParamString("merchant_name"),
		WrtIDs:       ctx.GetParamArrayString("wrt"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.SPVGetSalesOrderList(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVWidget(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppSPVWidgetRequest{
		SiteIdGp: ctx.Request().Context().Value(constants.KeySiteID).(string),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.SPVWidget(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVGetSalesOrderDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	sopNumber := ctx.GetParamString("sop_number")

	ctx.ResponseData, err = h.ServiceHelperApp.SPVGetSalesOrderDetail(ctx.Request().Context(), sopNumber)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVRejectSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppSopNumberRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.SpvId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.SPVRejectSalesOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVAcceptSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppSopNumberRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.SpvId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.SPVAcceptSalesOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVGetWrtMonitoring(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppGetWrtMonitoringRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.SiteId = ctx.Request().Context().Value(constants.KeySiteID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.SPVGetWrtMonitoring(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) SPVGetWrtMonitoringDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppWrtMonitoringDetailRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.SiteId = ctx.Request().Context().Value(constants.KeySiteID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.SPVGetWrtMonitoringDetail(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerGetSalesOrderList(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var req dto.HelperAppGetSalesOrderToCheckRequest
	req = dto.HelperAppGetSalesOrderToCheckRequest{
		Offset:       page.Offset,
		Limit:        page.Limit,
		SiteId:       ctx.Request().Context().Value(constants.KeySiteID).(string),
		SopNumber:    ctx.GetParamString("sop_number"),
		Statuses:     ctx.GetParamArrayInt("status"),
		WrtIDs:       ctx.GetParamArrayString("wrt"),
		CustomerName: ctx.GetParamString("merchant_name"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerGetSalesOrderList(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerWidget(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppCheckerWidgetRequest{
		CheckerId: ctx.Request().Context().Value(constants.KeyUserID).(string),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerWidget(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerGetSalesOrderDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppCheckerGetSalesOrderDetailRequest{
		SopNumber: ctx.GetParamString("sop_number"),
		CheckerId: ctx.Request().Context().Value(constants.KeyUserID).(string),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerGetSalesOrderDetail(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerStartChecking(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppCheckerStartCheckingRequest

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.CheckerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerStartChecking(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerSubmitChecking(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppCheckerSubmitCheckingRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.CheckerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerSubmitChecking(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerRejectSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppCheckerRejectSalesOrderRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.CheckerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerRejectSalesOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerGetDeliveryKoli(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	sopNumber := ctx.GetParamString("sop_number")

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerGetDeliveryKoli(ctx.Request().Context(), sopNumber)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerAcceptSalesOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.HelperAppCheckerAcceptSalesOrderRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	req.CheckerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerAcceptSalesOrder(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerHistory(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.HelperAppCheckerHistoryRequest{
		Offset:       page.Offset,
		Limit:        page.Limit,
		SopNumber:    ctx.GetParamString("sop_number"),
		WrtIdGP:      ctx.GetParamString("wrt_id_gp"),
		CheckerId:    ctx.Request().Context().Value(constants.KeyUserID).(string),
		CustomerName: ctx.GetParamString("merchant_name"),
	}

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerHistory(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h HelperAppHandler) CheckerHistoryDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	req := dto.HelperAppCheckerHistoryDetailRequest{
		SopNumber: ctx.GetParamString("sop_number"),
	}

	req.CheckerId = ctx.Request().Context().Value(constants.KeyUserID).(string)

	ctx.ResponseData, err = h.ServiceHelperApp.CheckerHistoryDetail(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
