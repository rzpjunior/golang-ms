package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ArchetypeHandler struct {
	Option           global.HandlerOptions
	ServiceArchetype service.IArchetypeService
}

// URLMapping declare endpoint with handler function.
func (h *ArchetypeHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceArchetype = service.NewArchetypeService()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized("public"))
}

func (h ArchetypeHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var total int64
	customerTypeId := ctx.GetParamInt("customer_type_id")

	var Archetypes []dto.ArchetypeResponse
	Archetypes, total, _ = h.ServiceArchetype.Get(ctx.Request().Context(), 0, 0, "", customerTypeId, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(Archetypes, total, page)

	return ctx.Serve(err)
}
