package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	Option         global.HandlerOptions
	ServicePayment service.IPaymentService
}

// URLMapping declare endpoint with handler function.
func (h *PaymentHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePayment = service.NewPaymentService()

	cMiddleware := middleware.NewMiddleware()

	r.POST("/list-payment", h.GetPaymentMethod, cMiddleware.Authorized("private"))
	r.POST("/xendit/invoice", h.getSalesInvoiceExternalXendit)
}

func (h PaymentHandler) GetPaymentMethod(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	req := &dto.PaymentMethodRequestGet{}
	if req.Session, err = service.CustomerSession(ctx); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.ResponseData, _ = h.ServicePayment.GetPaymentMethod(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	// ctx.DataList(PaymentMethods, total, page)
	//ctx.ResponseData = Banks

	return ctx.Serve(err)
}

func (h PaymentHandler) getSalesInvoiceExternalXendit(c echo.Context) (e error) {
	ctx := c.(*edenlabs.Context)
	var req *dto.PaymentRequest

	if e = ctx.Bind(&req); e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}
	// if req.Session, e = service.CustomerSession(ctx); e != nil {
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
	// 	return ctx.Serve(e)
	// }

	ctx.ResponseData, e = h.ServicePayment.GetInvoiceXendit(ctx.Request().Context(), req)
	if e != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, e).Print()
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
