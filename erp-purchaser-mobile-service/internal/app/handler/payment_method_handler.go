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

type PaymentMethodHandler struct {
	Option               global.HandlerOptions
	ServicePaymentMethod service.IPaymentMethodService
}

// URLMapping declare endpoint with handler function.
func (h *PaymentMethodHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePaymentMethod = service.NewPaymentMethodService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
}

func (h PaymentMethodHandler) Get(c echo.Context) (err error) {
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

	var paymentMethods []*dto.PaymentMethodResponse
	var total int64
	paymentMethods, err = h.ServicePaymentMethod.Get(ctx.Request().Context(), dto.PaymentMethodListRequest{
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

	ctx.DataList(paymentMethods, total, page)

	return ctx.Serve(err)
}

func (h PaymentMethodHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var paymentMethod *dto.PaymentMethodResponse

	id := c.Param("id")

	paymentMethod, err = h.ServicePaymentMethod.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = paymentMethod

	return ctx.Serve(err)
}
