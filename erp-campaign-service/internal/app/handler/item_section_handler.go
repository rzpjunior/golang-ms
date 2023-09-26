package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ItemSectionHandler struct {
	Option              global.HandlerOptions
	ServicesItemSection service.IItemSectionService
}

// URLMapping implements router.RouteHandlers
func (h *ItemSectionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItemSection = service.NewItemSectionService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("isc_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("isc_rdd"))
	r.POST("", h.Create, cMiddleware.Authorized("isc_crt"))
	r.PUT("/:id", h.Update, cMiddleware.Authorized("isc_upd"))
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized("isc_arc"))
}

func (h ItemSectionHandler) Get(c echo.Context) (err error) {
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
	statuses := ctx.GetParamArrayInt("status")
	regionID := ctx.GetParamString("region_id")
	archetypeID := ctx.GetParamString("archetype_id")
	sectionType := ctx.GetParamInt("type")
	param := &dto.ItemSectionRequestGet{
		Offset:      int64(page.Start),
		Limit:       int64(page.Limit),
		RegionID:    regionID,
		ArchetypeID: archetypeID,
		Search:      search,
		OrderBy:     orderBy,
		Type:        int8(sectionType),
	}
	for _, v := range statuses {
		param.Status = append(param.Status, int32(v))
	}

	var itemSections []*dto.ItemSectionResponse
	var total int64
	itemSections, total, err = h.ServicesItemSection.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(itemSections, total, page)

	return ctx.Serve(err)
}

func (h ItemSectionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var itemSection dto.ItemSectionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	itemSection, err = h.ServicesItemSection.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = itemSection

	return ctx.Serve(err)
}

func (h ItemSectionHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ItemSectionRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemSection.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemSectionHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ItemSectionRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemSection.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemSectionHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ItemSectionRequestArchive
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemSection.Archive(ctx.Request().Context(), id, req)

	return ctx.Serve(err)
}
