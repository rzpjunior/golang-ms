package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type MembershipHandler struct {
	Option             global.HandlerOptions
	ServicesMembership service.IMembershipService
}

// URLMapping implements router.RouteHandlers
func (h *MembershipHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesMembership = service.NewMembershipService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("/level", h.GetMembershipLevelList, cMiddleware.Authorized())
	r.GET("/level/:id", h.GetMembeshipLevelDetail, cMiddleware.Authorized())
	r.GET("/checkpoint", h.GetMembershipCheckpointList, cMiddleware.Authorized())
	r.GET("/checkpoint/:id", h.GetMembershipCheckpointDetail, cMiddleware.Authorized())
}

func (h MembershipHandler) GetMembershipLevelList(c echo.Context) (err error) {
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

	req := &dto.MembershipLevelRequestGet{
		Search:  search,
		OrderBy: orderBy,
		Status:  int8(status),
		Offset:  int64(page.Start),
		Limit:   int64(page.Limit),
	}

	var membershipLevels []*dto.MembershipLevelResponse
	var total int64
	membershipLevels, total, err = h.ServicesMembership.GetMembershipLevelList(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(membershipLevels, total, page)

	return ctx.Serve(err)
}

func (h MembershipHandler) GetMembeshipLevelDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var membershipLevel dto.MembershipLevelResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	membershipLevel, err = h.ServicesMembership.GetMembeshipLevelDetail(ctx.Request().Context(), id, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = membershipLevel

	return ctx.Serve(err)
}

func (h MembershipHandler) GetMembershipCheckpointList(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	// get params filter
	membershipLevelID := ctx.GetParamInt("level")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")

	req := &dto.MembershipCheckpointRequestGet{
		MembershipLevelID: int64(membershipLevelID),
		OrderBy:           orderBy,
		Status:            int8(status),
		Offset:            int64(page.Start),
		Limit:             int64(page.Limit),
	}

	var membershipCheckpoints []*dto.MembershipCheckpointResponse
	var total int64
	membershipCheckpoints, total, err = h.ServicesMembership.GetMembershipCheckpointList(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(membershipCheckpoints, total, page)

	return ctx.Serve(err)
}

func (h MembershipHandler) GetMembershipCheckpointDetail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var membershipCheckpoint dto.MembershipCheckpointResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	membershipCheckpoint, err = h.ServicesMembership.GetMembershipCheckpointDetail(ctx.Request().Context(), id, 0)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = membershipCheckpoint

	return ctx.Serve(err)
}
