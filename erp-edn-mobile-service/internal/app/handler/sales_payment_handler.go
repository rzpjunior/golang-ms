package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesPaymentHandler struct {
	Option               global.HandlerOptions
	ServicesSalesPayment service.ISalesPaymentService
	ServicesAuth         service.IAuthService
}

// URLMapping implements router.RouteHandlers
func (h *SalesPaymentHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesPayment = service.NewServiceSalesPayment()
	h.ServicesAuth = service.NewServiceAuth()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailById, cMiddleware.Authorized("edn_app"))
	r.POST("", h.CreateGP, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGP, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/:id", h.DetailGP, cMiddleware.Authorized("edn_app"))
	r.GET("/performance", h.PerformancePayment, cMiddleware.Authorized("edn_app"))
}

func (h SalesPaymentHandler) Get(c echo.Context) (err error) {
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

	var datas []*dto.SalesPaymentResponse
	var total int64
	datas, total, err = h.ServicesSalesPayment.Get(ctx.Request().Context(), dto.SalesPaymentListRequest{
		Limit:               int32(limit),
		Offset:              int32(offset),
		Status:              status,
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

func (h SalesPaymentHandler) DetailById(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var si *dto.SalesPaymentResponse
	var id string
	id = ctx.Param("id")

	si, err = h.ServicesSalesPayment.GetByID(ctx.Request().Context(), dto.SalesPaymentDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = si

	return ctx.Serve(err)
}

func (h SalesPaymentHandler) CreateGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// var session dto.UserResponse
	var po *dto.SalesPaymentResponse

	var req dto.CreateSalesPaymentRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	// if session, err = h.ServicesAuth.Session(ctx.Request().Context()); err != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return ctx.Serve(err)
	// }
	// req.RegionID = session.Region.Code
	po, err = h.ServicesSalesPayment.CreateGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h SalesPaymentHandler) GetGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	docnumbr := ctx.GetParamString("docnumbr")
	custnmbr := ctx.GetParamString("custnmbr")
	sopnumbe := ctx.GetParamString("sopnumbe")
	gnlRegion := ctx.GetParamString("gnl_region")
	locncode := ctx.GetParamString("locncode")
	docDateFrom := ctx.GetParamString("docdate_from")
	docDateto := ctx.GetParamString("docdate_to")
	siDocdateFrom := ctx.GetParamString("si_docdate_from")
	siDocdateTo := ctx.GetParamString("si_docdate_to")
	soDocdateFrom := ctx.GetParamString("so_docdate_from")
	soDocdateTo := ctx.GetParamString("so_docdate_to")

	var po []*dto.SalesPaymentGP
	var total int64
	po, total, err = h.ServicesSalesPayment.GetListGP(ctx.Request().Context(), dto.GetSalesPaymentGPListRequest{
		Limit:         int32(page.Limit),
		Offset:        int32(page.Offset),
		Docnumbr:      docnumbr,
		DocdateFrom:   docDateFrom,
		DocdateTo:     docDateto,
		SiDocdateFrom: siDocdateFrom,
		SiDocdateTo:   siDocdateTo,
		SoDocdateFrom: soDocdateFrom,
		SoDocdateTo:   soDocdateTo,
		Custnmbr:      custnmbr,
		Sopnumbe:      sopnumbe,
		GnlRegion:     gnlRegion,
		Locncode:      locncode,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(po, total, page)

	return ctx.Serve(err)
}

func (h SalesPaymentHandler) DetailGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.SalesPaymentGP

	var id string
	id = ctx.Param("id")

	po, err = h.ServicesSalesPayment.GetDetailGP(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}

func (h SalesPaymentHandler) PerformancePayment(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

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

	_, summaryPaymentPerformance, err := h.ServicesSalesPayment.GetPaymentPerformance(ctx.Request().Context(), dto.PerformancePaymentRequest{
		Limit:               int32(limit),
		Offset:              int32(offset),
		Status:              status,
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

	ctx.ResponseData = summaryPaymentPerformance

	return ctx.Serve(err)
}
