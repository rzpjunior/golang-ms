package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type DivisionHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *DivisionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
	r.GET("/:id", h.Detail, middleware.NewMiddleware().Authorized())
	r.POST("", h.Create, middleware.NewMiddleware().Authorized("dvs_crt"))
	r.PUT("/:id", h.Update, middleware.NewMiddleware().Authorized("dvs_upd"))
	r.PUT("/archive/:id", h.Archive, middleware.NewMiddleware().Authorized("dvs_arc"))
	r.GET("/default", h.GetByCustomerType, middleware.NewMiddleware().Authorized())
}

func (h DivisionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")

	var divisions []*dto.DivisionResponse
	var total int64
	divisions, total, err = sDivision.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(divisions, total, page)

	return ctx.Serve(err)
}

func (h DivisionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	var division dto.DivisionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	division, err = sDivision.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = division

	return ctx.Serve(err)
}

func (h DivisionHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	var req dto.DivisionRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sDivision.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h DivisionHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	var req dto.DivisionRequestUpdate

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sDivision.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h DivisionHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = sDivision.Archive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h DivisionHandler) GetByCustomerType(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sDivision := service.ServiceDivision()

	customerTypeID := ctx.GetParamString("customer_type_id")

	var division *dto.DivisionResponse
	division, err = sDivision.GetDivisonByCustomerType(ctx.Request().Context(), customerTypeID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.Data(division)
	return ctx.Serve(err)
}
