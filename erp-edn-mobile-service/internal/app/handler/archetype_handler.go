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

type ArchetypeHandler struct {
	Option            global.HandlerOptions
	ServicesArchetype service.IArchetypeService
}

// URLMapping implements router.RouteHandlers
func (h *ArchetypeHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesArchetype = service.NewServiceArchetype()

	cMiddleware := middleware.NewMiddleware()
	// GP Integrated
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h ArchetypeHandler) GetGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	gnlArchetypeId := ctx.GetParamString("gnl_archetype_id")
	gnlArchetypedescription := ctx.GetParamString("gnl_archetypedescription")
	gnlCustTypeId := ctx.GetParamString("gnl_cust_type_id")
	inactive := ctx.GetParamString("inactive")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var archerytpes []*dto.ArchetypeGP
	var total int64
	archerytpes, total, err = h.ServicesArchetype.GetGP(ctx.Request().Context(), dto.GetArchetypeGPListRequest{
		Limit:                   int32(limit),
		Offset:                  int32(offset),
		GnlArchetypeId:          gnlArchetypeId,
		GnlArchetypedescription: gnlArchetypedescription,
		GnlCustTypeId:           gnlCustTypeId,
		Inactive:                inactive,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(archerytpes, total, page)

	return ctx.Serve(err)
}

func (h ArchetypeHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var archetype *dto.ArchetypeGP

	var id string
	id = ctx.Param("id")

	archetype, err = h.ServicesArchetype.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = archetype

	return ctx.Serve(err)
}
