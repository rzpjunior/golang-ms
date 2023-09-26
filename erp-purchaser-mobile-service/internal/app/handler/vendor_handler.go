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

type VendorHandler struct {
	Option        global.HandlerOptions
	ServiceVendor service.IVendorService
}

// URLMapping declare endpoint with handler function.
func (h *VendorHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceVendor = service.NewVendorService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
	r.POST("", h.Create, cMiddleware.Authorized("purchaser_app"))
}

func (h VendorHandler) Get(c echo.Context) (err error) {
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
	vendorOrg := ctx.GetParamString("VendorOrg")
	state := ctx.GetParamString("state")
	var vendors []*dto.VendorResponse
	var total int64
	vendors, total, err = h.ServiceVendor.Get(ctx.Request().Context(), &dto.VendorListRequest{
		Limit:     int32(page.Limit),
		Offset:    int32(page.Offset),
		Status:    int32(status),
		Search:    search,
		OrderBy:   orderBy,
		VendorOrg: vendorOrg,
		State:     state,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(vendors, total, page)

	return ctx.Serve(err)
}

func (h VendorHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var vendor *dto.VendorResponse

	id := c.Param("id")

	vendor, err = h.ServiceVendor.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = vendor

	return ctx.Serve(err)
}

func (h VendorHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.VendorRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var res *dto.VendorRequestCreateResponse
	res, err = h.ServiceVendor.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}
