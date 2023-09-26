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

type CheckbookHandler struct {
	Option            global.HandlerOptions
	ServicesCheckbook service.ICheckbookService
}

// URLMapping implements router.RouteHandlers
func (h *CheckbookHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesCheckbook = service.NewServiceCheckbook()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.GetGp, cMiddleware.Authorized("edn_app"))
}

func (h CheckbookHandler) GetGp(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	regionID := ctx.GetParamString("region_id")

	var pm []*dto.CheckbookGP
	var total int64
	pm, total, err = h.ServicesCheckbook.GetGP(ctx.Request().Context(), dto.CheckbookListRequest{
		// Limit:  int32(page.Limit),
		// Offset: int32(page.Offset),
		// Search: search,
		RegionID: regionID,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(pm, total, page)

	return ctx.Serve(err)
}

// func (h CheckbookHandler) DetailGp(c echo.Context) (err error) {
// 	ctx := c.(*edenlabs.Context)

// 	var pm *dto.CheckbookGP

// 	var id string
// 	id = ctx.GetParamString("id")

// 	pm, err = h.ServicesCheckbook.GetDetaiGPlById(ctx.Request().Context(), id)
// 	if err != nil {
// 		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
// 		return ctx.Serve(err)
// 	}

// 	ctx.ResponseData = pm

// 	return ctx.Serve(err)
// }
