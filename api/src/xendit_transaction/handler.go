package xendit_transaction

import (
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

// Handler collection handler for privilege.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("/fixedva", h.create)
	r.POST("/invoice/paid", h.invoicePaid)
	r.POST("/invoice/expired", h.invoiceExpired)
}

func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r fixedVaRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = Save(r)
	}
	return ctx.Serve(e)
}

func (h *Handler) invoicePaid(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r invoicePaidRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = SaveInvoice(r)
	}
	return ctx.Serve(e)
}

func (h *Handler) invoiceExpired(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r invoiceExpiredRequest

	if e = ctx.Bind(&r); e == nil {
		ctx.ResponseData, e = ExpiredInvoice(r)
	}
	return ctx.Serve(e)
}
