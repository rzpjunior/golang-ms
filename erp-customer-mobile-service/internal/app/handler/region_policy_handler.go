package handler

import (
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type RegionPolicyHandler struct {
	Option              global.HandlerOptions
	ServiceRegionPolicy service.IRegionPolicyService
}

// URLMapping declare endpoint with handler function.
func (h *RegionPolicyHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceRegionPolicy = service.NewRegionPolicyService()
	cMiddleware := middleware.NewMiddleware()

	r.POST("/region_policy", h.GetRegionPolicy, cMiddleware.Authorized("public"))
}

func (h RegionPolicyHandler) GetRegionPolicy(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)
	var req dto.RegionPolicyMobileRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	var total int64
	admDivId, _ := strconv.Atoi(req.Data.AdmDivisionID)

	var RegionPolicy dto.RegionPolicy
	RegionPolicy, total, _ = h.ServiceRegionPolicy.Get(ctx.Request().Context(), 0, 0, "", req.Data.AdmDivisionID, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	ctx.ResponseData = RegionPolicy

	fmt.Println(total, admDivId)
	return ctx.Serve(err)
}
