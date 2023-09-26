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

type PaymentTermHandler struct {
	Option             global.HandlerOptions
	ServicePaymentTerm service.IPaymentTermService
}

// URLMapping declare endpoint with handler function.
func (h *PaymentTermHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePaymentTerm = service.NewPaymentTermService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
}

func (h PaymentTermHandler) Get(c echo.Context) (err error) {
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

	var paymentTerms []*dto.PaymentTermResponse
	var total int64
	paymentTerms, err = h.ServicePaymentTerm.Get(ctx.Request().Context(), dto.PaymentTermListRequest{
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

	ctx.DataList(paymentTerms, total, page)

	return ctx.Serve(err)
}

func (h PaymentTermHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var paymentTerm *dto.PaymentTermResponse

	id := c.Param("id")

	paymentTerm, err = h.ServicePaymentTerm.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = paymentTerm

	return ctx.Serve(err)
}
