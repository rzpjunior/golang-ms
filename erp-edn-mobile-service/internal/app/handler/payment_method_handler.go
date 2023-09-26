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

type PaymentMethodHandler struct {
	Option                global.HandlerOptions
	ServicesPaymentMethod service.IPaymentMethodService
}

// URLMapping implements router.RouteHandlers
func (h *PaymentMethodHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPaymentMethod = service.NewServicePaymentMethod()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
	r.GET("/:id", h.DetailGp, cMiddleware.Authorized("edn_app"))
}

func (h PaymentMethodHandler) GetGp(c echo.Context) (err error) {
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

	var pm []*dto.PaymentMethodGP
	var total int64
	pm, total, err = h.ServicesPaymentMethod.Get(ctx.Request().Context(), dto.PaymentMethodListRequest{
		// Limit:  int32(page.Limit),
		// Offset: int32(page.Offset),
		Search: search,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(pm, total, page)

	return ctx.Serve(err)
}

func (h PaymentMethodHandler) DetailGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var pm *dto.PaymentMethodGP

	var id string
	id = ctx.GetParamString("id")

	pm, err = h.ServicesPaymentMethod.GetDetaiGPlById(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = pm

	return ctx.Serve(err)
}
