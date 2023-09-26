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
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
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
	search := ctx.GetParamString("search")
	status := ctx.GetParamInt("status")
	customerTypeID := ctx.GetParamString("customer_type_id")

	param := &dto.ArchetypeGetListRequest{
		Offset:         page.Start,
		Limit:          page.Limit,
		Search:         search,
		Status:         int8(status),
		CustomerTypeID: customerTypeID,
	}

	var archetypes []*dto.ArchetypeResponse
	var total int64
	archetypes, total, err = h.ServicesArchetype.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(archetypes, total, page)

	return ctx.Serve(err)
}

func (h ArchetypeHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesArchetype.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
