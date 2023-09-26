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

type VendorOrganizationHandler struct {
	Option                    global.HandlerOptions
	ServiceVendorOrganization service.IVendorOrganizationService
}

// URLMapping declare endpoint with handler function.
func (h *VendorOrganizationHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceVendorOrganization = service.NewVendorOrganizationService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
}

func (h VendorOrganizationHandler) Get(c echo.Context) (err error) {
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

	var vendors []*dto.VendorOrganizationResponse
	var total int64
	vendors, total, err = h.ServiceVendorOrganization.Get(ctx.Request().Context(), &dto.VendorOrganizationListRequest{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Offset),
		Status:  int32(status),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(vendors, total, page)

	return ctx.Serve(err)
}

func (h VendorOrganizationHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var vendor *dto.VendorOrganizationResponse

	id := c.Param("id")

	vendor, err = h.ServiceVendorOrganization.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = vendor

	return ctx.Serve(err)
}
