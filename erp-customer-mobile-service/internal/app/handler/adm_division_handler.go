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

type AdmDivisionHandler struct {
	Option             global.HandlerOptions
	ServiceAdmDivision service.IAdmDivisionService
}

// URLMapping declare endpoint with handler function.
func (h *AdmDivisionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceAdmDivision = service.NewAdmDivisionService()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized("public"))
	r.POST("/gp", h.GetGP, cMiddleware.Authorized("public"))
	r.POST("/search", h.Search, cMiddleware.Authorized("public"))
}

func (h AdmDivisionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var total int64
	subDistrictID := ctx.GetParamInt("sub_district_id")
	search := ctx.GetParamString("search")
	regionID := ctx.GetParamInt("region_id")
	provinceID := ctx.GetParamInt("province_id")
	cityID := ctx.GetParamInt("city_id")
	districtID := ctx.GetParamInt("district_id")

	var AdmDivisions []dto.AdmDivisionResponse
	AdmDivisions, total, _ = h.ServiceAdmDivision.Get(ctx.Request().Context(), page.Offset, page.Limit, search, subDistrictID, 0, regionID, provinceID, cityID, districtID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisions, total, page)
	//ctx.ResponseData = AdmDivisions

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) GetGP(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		total        int64
		AdmDivisions []*dto.AdmDivisionGPResponse
		req          dto.GetAdmDivisionRequest
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	AdmDivisions, total, _ = h.ServiceAdmDivision.GetGP(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisions, total, &edenlabs.Paginator{})

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) Search(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		total        int64
		AdmDivisions []*dto.AdmDivisionGPResponse
		req          dto.SearchAdmDivisionRequest
	)

	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	AdmDivisions, total, _ = h.ServiceAdmDivision.Search(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisions, total, &edenlabs.Paginator{})

	return ctx.Serve(err)
}
