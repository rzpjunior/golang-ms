package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type MenuHandler struct {
	Option global.HandlerOptions
}

// URLMapping declare endpoint with handler function.
func (h *MenuHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	r.GET("", h.Get, middleware.NewMiddleware().Authorized())
}

func (h MenuHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	sMenu := service.ServiceMenu()

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	status := ctx.GetParamInt("status")
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")

	var menus []*dto.MenuResponse
	var total int64
	menus, total, err = sMenu.Get(ctx.Request().Context(), page.Start, page.Limit, status, search, orderBy)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(menus, total, page)

	return ctx.Serve(err)
}
