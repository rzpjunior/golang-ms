package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type BannerHandler struct {
	Option         global.HandlerOptions
	ServicesBanner service.IBannerService
}

// URLMapping implements router.RouteHandlers
func (h *BannerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesBanner = service.NewBannerService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("bnr_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("bnr_rdd"))
	r.POST("", h.Create, cMiddleware.Authorized("bnr_crt"))
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized("bnr_arc"))
}

func (h BannerHandler) Get(c echo.Context) (err error) {
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
	regionID := ctx.GetParamString("region_id")
	archetypeID := ctx.GetParamString("archetype_id")

	var banners []*dto.BannerResponse
	var total int64
	param := &dto.BannerRequestGet{
		Offset:      int64(page.Start),
		Limit:       int64(page.Limit),
		RegionID:    regionID,
		ArchetypeID: archetypeID,
		Search:      search,
		OrderBy:     orderBy,
		Status:      []int32{int32(status)},
	}

	banners, total, err = h.ServicesBanner.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(banners, total, page)

	return ctx.Serve(err)
}

func (h BannerHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var banner dto.BannerResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	banner, err = h.ServicesBanner.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = banner

	return ctx.Serve(err)
}

func (h BannerHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.BannerRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesBanner.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h BannerHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.BannerRequestArchive
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesBanner.Archive(ctx.Request().Context(), id, req)

	return ctx.Serve(err)
}
