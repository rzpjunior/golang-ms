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

type ItemHandler struct {
	Option       global.HandlerOptions
	ServicesItem service.IItemService
	ServiceAuth  service.IAuthService
}

// URLMapping implements router.RouteHandlers
func (h *ItemHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItem = service.NewServiceItem()
	h.ServiceAuth = service.NewServiceAuth()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
	r.GET("/gt", h.GetProductGT, cMiddleware.Authorized("edn_app"))
	// r.GET("/gt", h.GetProductGT)
	r.GET("/:id", h.DetailById, cMiddleware.Authorized("edn_app"))
}

func (h ItemHandler) Get(c echo.Context) (err error) {
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
	customerID := ctx.GetParamString("customer_id")

	var items []*dto.ItemResponse
	var total int64
	items, err = h.ServicesItem.Get(ctx.Request().Context(), dto.ItemListRequest{
		Limit:      int32(limit),
		Offset:     int32(offset),
		Status:     int32(status),
		Search:     search,
		OrderBy:    orderBy,
		CustomerID: customerID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h ItemHandler) GetProductGT(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	var session dto.UserResponse
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	if session, err = h.ServiceAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")

	// WRH0007
	var items []*dto.ItemResponse
	var total int64
	items, err = h.ServicesItem.GetProductGT(ctx.Request().Context(), dto.ItemListRequest{
		Limit:   int32(page.PerPage),
		Offset:  int32(page.Page - 1),
		Status:  int32(status),
		SiteID:  session.SiteID,
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h ItemHandler) DetailById(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var item *dto.ItemResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	item, err = h.ServicesItem.GetByID(ctx.Request().Context(), dto.ItemDetailRequest{
		Id: int32(id),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = item

	return ctx.Serve(err)
}
