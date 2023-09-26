package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PaymentTermHandler struct {
	Option              global.HandlerOptions
	ServicesPaymentTerm service.IPaymentTermService
}

// URLMapping implements router.RouteHandlers
func (h *PaymentTermHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPaymentTerm = service.NewServicePaymentTerm()

	cMiddleware := middleware.NewMiddleware()

	r.GET("", h.Get, cMiddleware.Authorized())
}

func (h PaymentTermHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	paymentUseFor := ctx.GetParamString("payment_usefor")

	req := &dto.GetPaymentTermRequest{
		Limit:         page.Limit,
		Offset:        page.Offset,
		PaymentUseFor: paymentUseFor,
	}

	paymentTerm, total, err := h.ServicesPaymentTerm.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.DataList(paymentTerm, total, page)

	return ctx.Serve(err)
}
