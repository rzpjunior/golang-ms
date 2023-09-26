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

type PaymentTermHandler struct {
	Option              global.HandlerOptions
	ServicesPaymentTerm service.IPaymentTermService
}

// URLMapping implements router.RouteHandlers
func (h *PaymentTermHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesPaymentTerm = service.NewServicePaymentTerm()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
}

func (h PaymentTermHandler) GetGp(c echo.Context) (err error) {
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
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var pm []*dto.PaymentTermGP
	//SalesPaymentTermGPResponse
	var total int64
	pm, total, err = h.ServicesPaymentTerm.GetGP(ctx.Request().Context(), dto.PaymentTermListRequest{
		Search: search,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(pm, total, page)

	return ctx.Serve(err)
}
