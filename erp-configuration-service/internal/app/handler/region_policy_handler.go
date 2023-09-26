package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type RegionPolicyHandler struct {
	Option              global.HandlerOptions
	ServiceRegionPolicy service.IRegionPolicyService
}

// URLMapping implements router.RouteHandlers
func (h *RegionPolicyHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceRegionPolicy = service.NewRegionPolicyService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
	r.PUT("/:id", h.Update, cMiddleware.Authorized("rgp_upd"))
}

func (h RegionPolicyHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	regionID := ctx.GetParamString("region_id")

	var RegionPolicys []dto.RegionPolicyResponse
	var total int64
	RegionPolicys, total, err = h.ServiceRegionPolicy.Get(ctx.Request().Context(), page.Start, page.Limit, search, orderBy, regionID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(RegionPolicys, total, page)

	return ctx.Serve(err)
}

func (h RegionPolicyHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var regionPolicy dto.RegionPolicyResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	code := ctx.GetParamString("code")
	regionID := ctx.GetParamString("region_id")

	regionPolicy, err = h.ServiceRegionPolicy.GetDetail(ctx.Request().Context(), id, code, regionID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = regionPolicy

	return ctx.Serve(err)
}

func (h RegionPolicyHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.RegionPolicyRequestUpdate

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

	ctx.ResponseData, err = h.ServiceRegionPolicy.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
