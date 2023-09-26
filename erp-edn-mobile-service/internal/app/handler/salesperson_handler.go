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

type SalespersonHandler struct {
	Option              global.HandlerOptions
	ServicesSalesperson service.ISalespersonService
}

// URLMapping implements router.RouteHandlers
func (h *SalespersonHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesperson = service.NewServiceSalesperson()

	cMiddleware := middleware.NewMiddleware()
	// r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	// r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h SalespersonHandler) Get(c echo.Context) (err error) {
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
	status := ctx.GetParamInt("status")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var sp []*dto.SalespersonResponse
	var total int64
	sp, err = h.ServicesSalesperson.Get(ctx.Request().Context(), dto.SalespersonListRequest{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(sp, total, page)

	return ctx.Serve(err)
}

func (h SalespersonHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var sp *dto.SalespersonResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	sp, err = h.ServicesSalesperson.GetDetailById(ctx.Request().Context(), dto.SalespersonDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = sp

	return ctx.Serve(err)
}

func (h SalespersonHandler) GetGp(c echo.Context) (err error) {
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
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")
	status := ctx.GetParamInt("status")

	var sp []*dto.SalesPerson
	var total int64
	sp, total, err = h.ServicesSalesperson.GetGP(ctx.Request().Context(), dto.SalespersonListRequest{
		Limit:       int32(limit),
		Offset:      int32(offset),
		OffsetQuery: int32(page.Offset),
		Search:      search,
		Status:      int32(status),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(sp, total, page)

	return ctx.Serve(err)
}

func (h SalespersonHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var sp *dto.SalesPerson

	var id string
	id = ctx.Param("id")

	sp, err = h.ServicesSalesperson.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = sp

	return ctx.Serve(err)
}
