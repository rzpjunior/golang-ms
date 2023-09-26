package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type SalesAssignmentObjectiveHandler struct {
	Option                           global.HandlerOptions
	ServicesSalesAssignmentObjective service.ISalesAssignmentObjectiveService
}

// URLMapping implements router.RouteHandlers
func (h *SalesAssignmentObjectiveHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesAssignmentObjective = service.NewSalesAssignmentObjectiveService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
	r.POST("", h.Create, cMiddleware.Authorized())
	r.PUT("/:id", h.Update, cMiddleware.Authorized())
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized())
	r.PUT("/unarchive/:id", h.UnArchive, cMiddleware.Authorized())
}

func (h SalesAssignmentObjectiveHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	codes := ctx.GetParamArrayString("codes")
	orderBy := ctx.GetParamString("order_by")

	var salesAssignmentObjectives []*dto.SalesAssignmentObjectiveResponse
	var total int64
	salesAssignmentObjectives, total, err = h.ServicesSalesAssignmentObjective.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, codes, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(salesAssignmentObjectives, total, page)

	return ctx.Serve(err)
}

func (h SalesAssignmentObjectiveHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var salesAssignmentObjective dto.SalesAssignmentObjectiveResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	salesAssignmentObjective, err = h.ServicesSalesAssignmentObjective.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = salesAssignmentObjective

	return ctx.Serve(err)
}

func (h SalesAssignmentObjectiveHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.SalesAssignmentObjectiveRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignmentObjective.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h SalesAssignmentObjectiveHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.SalesAssignmentObjectiveRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignmentObjective.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h SalesAssignmentObjectiveHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignmentObjective.Archive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}

func (h SalesAssignmentObjectiveHandler) UnArchive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesSalesAssignmentObjective.UnArchive(ctx.Request().Context(), id)

	return ctx.Serve(err)
}
