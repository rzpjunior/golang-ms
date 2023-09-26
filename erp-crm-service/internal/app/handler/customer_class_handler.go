package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type CustomerClassHandler struct {
	Option                global.HandlerOptions
	ServicesCustomerClass service.ICustomerClassService
}

// URLMapping implements router.RouteHandlers
func (h *CustomerClassHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCustomerClass = service.NewCustomerClassService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/detail", h.GetDetail, cMiddleware.Authorized())
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
	search := ctx.GetParamString("search")

	param := &dto.CustomerClassGetListRequest{
		Offset: int64(page.Start),
		Limit:  int64(page.Limit),
		Search: search,
	}

	var customerClasss []*dto.CustomerClassResponse
	var total int64
	customerClasss, total, err = h.ServicesCustomerClass.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(customerClasss, total, page)

	return ctx.Serve(err)
}

func (h CustomerClassHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := ctx.GetParamString("id")

	ctx.ResponseData, err = h.ServicesCustomerClass.GetDetail(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	return ctx.Serve(err)
}
