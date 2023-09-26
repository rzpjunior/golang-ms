package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type VoucherHandler struct {
	Option          global.HandlerOptions
	ServicesVoucher service.IVoucherService
}

// URLMapping implements router.RouteHandlers
func (h *VoucherHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesVoucher = service.NewVoucherService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("vou_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("vou_rdd"))
	r.POST("", h.Create, cMiddleware.Authorized("vou_crt"))
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized("vou_arc"))
	r.POST("/bulky", h.CreateBulky, cMiddleware.Authorized("vou_blk_imp"))

}

func (h VoucherHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var (
		page  *edenlabs.Paginator
		req   *dto.VoucherRequestGet
		data  []*dto.VoucherResponse
		total int64
	)
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	regionID := ctx.GetParamString("region_id")
	archetypeID := ctx.GetParamString("archetype_id")
	voucherType := ctx.GetParamInt("type")
	membershipLevelID := ctx.GetParamInt("membership_level_id")
	membershipCheckpointID := ctx.GetParamInt("membership_checkpoint_id")
	CustomerID := ctx.GetParamInt("customer_id")

	req = &dto.VoucherRequestGet{
		Search:                 search,
		OrderBy:                orderBy,
		Status:                 int8(status),
		RegionID:               regionID,
		ArchetypeID:            archetypeID,
		Type:                   int8(voucherType),
		MembershipLevelID:      int64(membershipLevelID),
		MembershipCheckpointID: int64(membershipCheckpointID),
		CustomerID:             int64(CustomerID),
		Limit:                  int64(page.Limit),
		Offset:                 int64(page.Start),
	}

	data, total, err = h.ServicesVoucher.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(data, total, page)

	return ctx.Serve(err)
}

func (h VoucherHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesVoucher.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}

func (h VoucherHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.VoucherRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesVoucher.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h VoucherHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.ResponseData, err = h.ServicesVoucher.Archive(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h VoucherHandler) CreateBulky(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.VoucherRequestBulky
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	err = h.ServicesVoucher.CreateBulky(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
