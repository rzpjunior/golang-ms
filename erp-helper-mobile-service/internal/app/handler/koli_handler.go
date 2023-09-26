package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-helper-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type KoliHandler struct {
	Option       global.HandlerOptions
	ServicesKoli service.IKoliService
}

// URLMapping implements router.RouteHandlers
func (h *KoliHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesKoli = service.NewServiceKoli()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.AuthorizedHelperMobile())
	r.GET("/:id", h.GetDetail, cMiddleware.AuthorizedHelperMobile())
}

func (h KoliHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	req := dto.GetKoliRequest{
		Offset:  page.Offset,
		Limit:   page.Limit,
		OrderBy: ctx.GetParamString("order_by"),
		Status:  ctx.GetParamInt("status"),
	}

	koli, total, err := h.ServicesKoli.GetKoli(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(koli, total, page)

	return ctx.Serve(err)
}

func (h KoliHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesKoli.GetKoliDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
