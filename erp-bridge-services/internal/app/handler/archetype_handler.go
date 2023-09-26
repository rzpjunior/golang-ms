package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ArchetypeHandler struct {
	Option            global.HandlerOptions
	ServicesArchetype service.IArchetypeService
}

// URLMapping implements router.RouteHandlers
func (h *ArchetypeHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesArchetype = service.NewArchetypeService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h ArchetypeHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")

	CustomerTypeID := ctx.GetParamInt("customer_type_id")

	var archetypes []dto.ArchetypeResponse
	var total int64
	archetypes, total, err = h.ServicesArchetype.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, int64(CustomerTypeID))
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(archetypes, total, page)

	return ctx.Serve(err)
}

func (h ArchetypeHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var archetype dto.ArchetypeResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	archetype, err = h.ServicesArchetype.GetDetail(ctx.Request().Context(), id, "")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = archetype

	return ctx.Serve(err)
}
