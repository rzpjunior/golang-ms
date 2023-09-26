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

type ItemSectionHandler struct {
	Option       global.HandlerOptions
	ServicesItem service.IItemService
}

// URLMapping implements router.RouteHandlers
func (h *ItemSectionHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItem = service.NewItemService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("/item", h.Get, cMiddleware.Authorized())
}

func (h ItemSectionHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	items, total, err := h.ServicesItem.GetListItemComplex(ctx.Request().Context(), &dto.ItemRequestGet{
		Offset:         page.Start,
		Limit:          page.Limit,
		Status:         ctx.GetParamInt("status"),
		Search:         ctx.GetParamString("search"),
		OrderBy:        ctx.GetParamString("order_by"),
		UomID:          ctx.GetParamString("uom_id"),
		ItemCategoryID: int64(ctx.GetParamInt("item_category_id")),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var itemSections []dto.ItemSectionResponse
	for _, item := range items {
		itemSections = append(itemSections, dto.ItemSectionResponse{
			ID:          item.ID,
			Code:        item.Code,
			Description: item.Description,
		})
	}

	ctx.DataList(itemSections, total, page)

	return ctx.Serve(err)
}
