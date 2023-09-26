package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ItemClassHandler struct {
	Option            global.HandlerOptions
	ServicesItemClass service.IItemClassService
}

// URLMapping implements router.RouteHandlers
func (h *ItemClassHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesItemClass = service.NewServiceItemClass()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
}

func (h ItemClassHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	search := ctx.GetParamString("search")
	req := dto.GetItemClassRequest{
		Limit:  page.Limit,
		Offset: page.Offset,
		Search: search,
	}

	itemClass, total, err := h.ServicesItemClass.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(itemClass, total, page)

	return ctx.Serve(err)
}

func (h ItemClassHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesItemClass.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
