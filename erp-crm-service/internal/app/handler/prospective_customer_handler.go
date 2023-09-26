package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type ProspectiveCustomerHandler struct {
	Option                      global.HandlerOptions
	ServicesProspectiveCustomer service.IProspectiveCustomerService
}

// URLMapping implements router.RouteHandlers
func (h *ProspectiveCustomerHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesProspectiveCustomer = service.NewProspectiveCustomerService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
	r.PUT("/decline/:id", h.Decline, cMiddleware.Authorized())
	r.POST("", h.Create, cMiddleware.Authorized())
	r.POST("/upgrade", h.Upgrade, cMiddleware.Authorized())
}

func (h ProspectiveCustomerHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// get params filters
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	archetypeID := ctx.GetParamString("archetype_id")
	CustomerTypeID := ctx.GetParamString("customer_type_id")
	regionID := ctx.GetParamString("region_id")
	salespersonID := ctx.GetParamString("salesperson_id")
	requestedBy := ctx.GetParamString("requested_by")
	customerID := ctx.GetParamString("customer_id")
	param := &dto.ProspectiveCustomerGetRequest{
		Search:         search,
		Offset:         int64(page.Start),
		Limit:          int64(page.Limit),
		Status:         int8(status),
		ArchetypeID:    archetypeID,
		CustomerID:     customerID,
		CustomerTypeID: CustomerTypeID,
		SalesPersonID:  salespersonID,
		RequestBy:      requestedBy,
		RegionID:       regionID,
		OrderBy:        orderBy,
	}
	var territories []*dto.ProspectiveCustomerResponse
	var total int64
	territories, total, err = h.ServicesProspectiveCustomer.Get(ctx.Request().Context(), param)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(territories, total, page)

	return ctx.Serve(err)
}

func (h ProspectiveCustomerHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var prospectiveCustomer *dto.ProspectiveCustomerResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	prospectiveCustomer, err = h.ServicesProspectiveCustomer.GetDetail(ctx.Request().Context(), &dto.ProspectiveCustomerGetDetailRequest{ID: id})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = prospectiveCustomer

	return ctx.Serve(err)
}

func (h ProspectiveCustomerHandler) Decline(c echo.Context) (err error) {

	ctx := c.(*edenlabs.Context)

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req dto.ProspectiveCustomerDecineRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesProspectiveCustomer.Decline(ctx.Request().Context(), req, id)

	return ctx.Serve(err)
}

func (h ProspectiveCustomerHandler) Create(c echo.Context) (err error) {

	ctx := c.(*edenlabs.Context)

	var req dto.ProspectiveCustomerCreateRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesProspectiveCustomer.Create(ctx.Request().Context(), &req)

	return ctx.Serve(err)
}

func (h ProspectiveCustomerHandler) Upgrade(c echo.Context) (err error) {

	ctx := c.(*edenlabs.Context)

	var req dto.ProspectiveCustomerUpgradeRequest
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData, err = h.ServicesProspectiveCustomer.Upgrade(ctx.Request().Context(), &req)

	return ctx.Serve(err)
}
