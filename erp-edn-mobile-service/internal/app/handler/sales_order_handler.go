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

type SalesOrderHandler struct {
	Option             global.HandlerOptions
	ServicesSalesOrder service.ISalesOrderService
}

// URLMapping implements router.RouteHandlers
func (h *SalesOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesSalesOrder = service.NewServiceSalesOrder()

	cMiddleware := middleware.NewMiddleware()
	r.POST("", h.Create, cMiddleware.Authorized("edn_app"))
}

func (h SalesOrderHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var po *dto.SalesOrderResponse

	var req dto.CreateSalesOrderRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	po, err = h.ServicesSalesOrder.Create(ctx.Request().Context(), &req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = po

	return ctx.Serve(err)
}
