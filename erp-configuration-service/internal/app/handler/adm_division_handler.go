package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type AdmDivisionHandler struct {
	Option              global.HandlerOptions
	ServicesAdmDivision service.IAdmDivisionService
}

// URLMapping implements router.RouteHandlers
func (h *AdmDivisionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAdmDivision = service.NewAdmDivisionService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
}

func (h AdmDivisionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	region := ctx.GetParamString("region")
	regionSearch := ctx.GetParamString("region_search")
	province := ctx.GetParamString("province")
	provinceSearch := ctx.GetParamString("province_search")
	city := ctx.GetParamString("city")
	citySearch := ctx.GetParamString("city_search")
	district := ctx.GetParamString("district")
	districtSearch := ctx.GetParamString("district_search")
	subdistrict := ctx.GetParamString("subdistrict")
	subdistrictSearch := ctx.GetParamString("subdistrict_search")

	param := &dto.AdmDivisionGetRequest{
		Region:            region,
		RegionSearch:      regionSearch,
		Province:          province,
		ProvinceSearch:    provinceSearch,
		City:              city,
		CitySearch:        citySearch,
		District:          district,
		DistrictSearch:    districtSearch,
		SubDistrict:       subdistrict,
		SubDistrictSearch: subdistrictSearch,
		Limit:             int64(page.Limit),
		Offset:            int64(page.Start),
	}

	var AdmDivisiones []*dto.AdmDivisionResponse
	var total int64
	AdmDivisiones, total, err = h.ServicesAdmDivision.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisiones, total, page)

	return ctx.Serve(err)
}
