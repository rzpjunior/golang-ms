package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
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

	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
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
	subDistrictId := ctx.GetParamString("sub_district_id")

	var AdmDivisions []dto.AdmDivisionResponse
	AdmDivisions, total, _ = h.ServiceAdmDivision.Get(ctx.Request().Context(), 0, 0, "", subDistrictId, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(AdmDivisions, total, page)
	//ctx.ResponseData = AdmDivisions

	return ctx.Serve(err)
}

func (h AdmDivisionHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var admDivision *dto.AdmDivisionResponse

	id := c.Param("id")

	admDivision, err = h.ServiceAdmDivision.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = admDivision

	return ctx.Serve(err)
}
