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

type AdmDivisionHandler struct {
	Option              global.HandlerOptions
	ServicesAdmDivision service.IAdmDivisionService
}

// URLMapping implements router.RouteHandlers
func (h *AdmDivisionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesAdmDivision = service.NewServiceAdmDivision()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("edn_app"))
	r.GET("/gp", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/detail", h.DetailGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/coverage", h.GetCoverageGp, cMiddleware.Authorized("edn_app"))
	r.GET("/gp/coverage/:id", h.GetCoverageDetailGp, cMiddleware.Authorized("edn_app"))
	r.POST("", h.GetGPAdmDiv, cMiddleware.Authorized("edn_app"))
}

func (h AdmDivisionHandler) GetGPAdmDiv(c echo.Context) (err error) {
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

	AdmDivisions, total, _ = h.ServicesAdmDivision.GetGPAdmDiv(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisions, total, &edenlabs.Paginator{})

	return ctx.Serve(err)
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

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var admDiv []*dto.AdmDivisionResponse
	var total int64
	admDiv, err = h.ServicesAdmDivision.GetAdmDivisions(ctx.Request().Context(), dto.AdmDivisionListRequest{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(admDiv, total, page)

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var admDiv *dto.AdmDivisionResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	admDiv, err = h.ServicesAdmDivision.GetAdmDivisionDetailById(ctx.Request().Context(), dto.AdmDivisionDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = admDiv

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) GetGp(c echo.Context) (err error) {
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
	code := ctx.GetParamString("code")
	region := ctx.GetParamString("region")
	state := ctx.GetParamString("state")
	city := ctx.GetParamString("city")
	dsitrcit := ctx.GetParamString("dsitrcit")
	subdistrict := ctx.GetParamString("subdistrict")
	admType := ctx.GetParamString("type")

	var adm []*dto.AdmDivisionGP
	var total int64

	adm, total, err = h.ServicesAdmDivision.GetGP(ctx.Request().Context(), dto.AdmDivisionListRequest{
		Limit:       int32(page.Limit),
		Offset:      int32(page.Offset),
		Search:      search,
		Code:        code,
		Region:      region,
		State:       state,
		City:        city,
		District:    dsitrcit,
		Subdistrict: subdistrict,
		Type:        admType,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(adm, total, page)

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var (
		adm  []*dto.AdmDivisionGP
		page *edenlabs.Paginator
	)

	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	id := ctx.GetParamString("id")
	divType := ctx.GetParamString("type")

	adm, err = h.ServicesAdmDivision.GetDetaiGPlById(ctx.Request().Context(), id, divType, page.Limit, page.Offset)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = adm

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) GetCoverageGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	gnlAdministrativeCode := ctx.GetParamString("gnl_administrative_code")
	gnlProvince := ctx.GetParamString("gnl_province")
	gnlCity := ctx.GetParamString("gnl_city")
	gnlDistrict := ctx.GetParamString("gnl_district")
	gnlSubdistrict := ctx.GetParamString("gnl_subdistrict")
	locncode := ctx.GetParamString("locncode")

	var adm []*dto.AdmDivisionCoverageGP
	var total int64
	adm, total, err = h.ServicesAdmDivision.GetCoverageGP(ctx.Request().Context(), dto.AdmDivisionCoverageListRequest{
		Limit:                 int32(page.Limit),
		Offset:                int32(page.Offset),
		GnlAdministrativeCode: gnlAdministrativeCode,
		GnlProvince:           gnlProvince,
		GnlCity:               gnlCity,
		GnlDistrict:           gnlDistrict,
		GnlSubdistrict:        gnlSubdistrict,
		Locncode:              locncode,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(adm, total, page)

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) GetCoverageDetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var adm *dto.AdmDivisionCoverageGP

	var id string
	id = ctx.GetParamString("id")

	adm, err = h.ServicesAdmDivision.GetCoverageDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = adm

	return ctx.Serve(err)
}
