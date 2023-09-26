package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ItemCategoryHandler struct {
	Option               global.HandlerOptions
	ServicesItemCategory service.IItemCategoryService
}

// URLMapping implements router.RouteHandlers
func (h *ItemCategoryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItemCategory = service.NewItemCategoryService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("itm_ctg_rdl"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("itm_ctg_rdd"))
	r.PUT("/:id/image", h.UpdateImage, cMiddleware.Authorized("itm_ctg_upl_img"))
	r.PUT("/archive/:id", h.Archive, cMiddleware.Authorized("itm_ctg_arc"))
	r.PUT("/unarchive/:id", h.Unarchive, cMiddleware.Authorized("itm_ctg_urc"))
	r.POST("", h.Create, cMiddleware.Authorized("itm_ctg_crt"))
	r.PUT("/:id", h.Update, cMiddleware.Authorized("itm_ctg_upd"))
}

func (h ItemCategoryHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
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

	var items []dto.ItemCategoryResponse
	var total int64
	items, total, err = h.ServicesItemCategory.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy, regionID)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var item dto.ItemCategoryResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	item, err = h.ServicesItemCategory.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = item

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) UpdateImage(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ItemCategoryImageRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemCategory.UpdateImage(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) Archive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemCategory.Archive(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) Unarchive(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemCategory.Unarchive(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.ItemCategoryRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemCategory.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h ItemCategoryHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ItemCategoryRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesItemCategory.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
