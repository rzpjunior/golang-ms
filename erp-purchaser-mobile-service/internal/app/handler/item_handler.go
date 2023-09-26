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

type ItemHandler struct {
	Option      global.HandlerOptions
	ServiceItem service.IItemService
}

// URLMapping declare endpoint with handler function.
func (h *ItemHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceItem = service.NewItemService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
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

	var items []*dto.ItemResponse
	var total int64
	items, err = h.ServiceItem.Get(ctx.Request().Context(), dto.ItemListRequest{
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

	ctx.DataList(items, total, page)

	return ctx.Serve(err)
}

func (h ItemHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var item *dto.ItemResponse

	id := c.Param("id")

	item, err = h.ServiceItem.GetByID(ctx.Request().Context(), dto.ItemDetailRequest{
		Id: id,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = item

	return ctx.Serve(err)
}
