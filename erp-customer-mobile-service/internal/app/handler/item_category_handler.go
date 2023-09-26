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

type ItemCategoryHandler struct {
	Option              global.HandlerOptions
	ServiceItemCategory service.IItemCategoryService
}

// URLMapping declare endpoint with handler function.
func (h *ItemCategoryHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceItemCategory = service.NewItemCategoryService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("", h.Get, cMiddleware.Authorized("public"))
}

func (h ItemCategoryHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.ItemCategoryMobileRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServiceItemCategory.Get(ctx.Request().Context(), req)

	return ctx.Serve(err)
}
