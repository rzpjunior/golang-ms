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

type CustomerClassHandler struct {
	Option                global.HandlerOptions
	ServicesCustomerClass service.ICustomerClassService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerClassHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerClass = service.NewServiceCustomerClass()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("edn_app"))
}

func (h CustomerClassHandler) Get(c echo.Context) (err error) {
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
	offset := ctx.GetParamInt("page") - 1
	limit := ctx.GetParamInt("per_page")

	var customerClass []*dto.CustomerClassResponse
	var total int64
	customerClass, total, err = h.ServicesCustomerClass.GetCustomerClass(ctx.Request().Context(), dto.CustomerClassRequest{
		Limit:   int64(limit),
		Offset:  int64(offset),
		Search:  search,
		OrderBy: orderBy,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerClass, total, page)

	return ctx.Serve(err)
}
