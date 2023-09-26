package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/labstack/echo/v4"
)

type SalesInvoiceHandler struct {
	Option               global.HandlerOptions
	ServicesSalesInvoice service.ISalesInvoiceService
}

// URLMapping implements router.RouteHandlers
func (h *SalesInvoiceHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesInvoice = service.NewServiceSalesInvoice()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailById, cMiddleware.Authorized("edn_app"))
	r.POST("", h.CreateGP, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGP, cMiddleware.Authorized("edn_app"))
	r.GET("/performance", h.PerformanceOrder, cMiddleware.Authorized("edn_app"))
}

func (h SalesInvoiceHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamString("status")
	siteID := ctx.GetParamString("site_id")
	recognitionDateFrom := ctx.GetParamDate("recognition_date_from")
	recognitionDateTo := ctx.GetParamDate("recognition_date_to")
	customerID := ctx.GetParamString("customer_id")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var datas []*dto.SalesInvoiceResponse
	var total int64
	datas, total, err = h.ServicesSalesInvoice.Get(ctx.Request().Context(), dto.SalesInvoiceListRequest{
		Limit:               int32(limit),
		Offset:              int32(offset),
		Status:              utils.ToString(status),
		Search:              search,
		OrderBy:             orderBy,
		SiteID:              siteID,
		RecognitionDateFrom: recognitionDateFrom,
		RecognitionDateTo:   recognitionDateTo,
		CustomerID:          customerID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(datas, total, page)

	return ctx.Serve(err)
}

func (h SalesInvoiceHandler) DetailById(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var si *dto.SalesInvoiceResponse

	id := c.Param("id")

	si, err = h.ServicesSalesInvoice.GetByID(ctx.Request().Context(), dto.SalesInvoiceDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = si

	return ctx.Serve(err)
}

func (h SalesInvoiceHandler) CreateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		req dto.CreateSalesInvoiceRequest
		si  *dto.SalesInvoiceResponse
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	si, err = h.ServicesSalesInvoice.CreateGP(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = si

	return ctx.Serve(err)
}

func (h SalesInvoiceHandler) GetGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	sopnumbe := ctx.GetParamString("sopnumbe")

	var si []*bridgeService.SalesInvoiceGP
	var total int64
	si, total, err = h.ServicesSalesInvoice.GetListGP(ctx.Request().Context(), dto.GetSalesInvoiceGPRequest{
		Limit:    int32(page.Limit),
		Offset:   int32(page.Offset),
		Sopnumbe: sopnumbe,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(si, total, page)

	return ctx.Serve(err)
}

func (h SalesInvoiceHandler) PerformanceOrder(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	status := ctx.GetParamString("status")
	siteID := ctx.GetParamString("site_id")
	recognitionDateFrom := ctx.GetParamDate("recognition_date_from")
	recognitionDateTo := ctx.GetParamDate("recognition_date_to")
	customerID := ctx.GetParamString("customer_id")

	_, summaryOrderPerformance, err := h.ServicesSalesInvoice.GetOrderPerformance(ctx.Request().Context(), dto.SalesInvoiceListRequest{
		Limit:               int32(page.Limit),
		Offset:              int32(page.Offset),
		Status:              utils.ToString(status),
		RecognitionDateFrom: recognitionDateFrom,
		RecognitionDateTo:   recognitionDateTo,
		SiteID:              siteID,
		CustomerID:          customerID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = summaryOrderPerformance

	return ctx.Serve(err)
}
