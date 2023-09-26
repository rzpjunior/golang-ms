package handler

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PurchasePlanHandler struct {
	Option              global.HandlerOptions
	ServicePurchasePlan service.IPurchasePlanService
	ServicesAuth        service.IAuthService
}

// URLMapping declare endpoint with handler function.
func (h *PurchasePlanHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePurchasePlan = service.NewPurchasePlanService()
	h.ServicesAuth = service.NewAuthService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
	r.GET("/summary", h.GetSummary, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/assign/:id", h.Assign, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/assign/cancel/:id", h.CancelAssign, cMiddleware.Authorized("purchaser_app"))
}

func (h PurchasePlanHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var session dto.UserResponse
	// get params pagination
	var page *edenlabs.Paginator
	page, err = edenlabs.NewPaginator(ctx)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	if session, err = h.ServicesAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	// get params filter
	search := ctx.GetParamString("search")
	orderBy := ctx.GetParamString("order_by")
	status := ctx.GetParamInt("status")
	siteID := ctx.GetParamString("site_id")
	fieldPurchaser := ctx.GetParamString("field_purchaser")
	recognitionDateFrom := ctx.GetParamDate("recognition_date_from")
	recognitionDateTo := ctx.GetParamDate("recognition_date_to")
	purchasePlanDateFrom := ctx.GetParamDate("purchase_plan_date_from")
	purchasePlanDateTo := ctx.GetParamDate("purchase_plan_date_to")
	// session:= Session(ctx)
	var purchasePlans []*dto.PurchasePlanResponse
	var total int64

	fmt.Println(session)
	if session.Site.ID != "" {
		siteID = session.Site.ID
	}
	if session.MainRole == "Field Purchaser" {
		fieldPurchaser = session.EmployeeCode
	}

	purchasePlans, total, err = h.ServicePurchasePlan.Get(ctx.Request().Context(), &dto.PurchasePlanListRequest{
		Limit:                int32(page.Limit),
		Offset:               int32(page.Offset),
		Status:               int32(status),
		Search:               search,
		OrderBy:              orderBy,
		SiteID:               siteID,
		RecognitionDateFrom:  recognitionDateFrom,
		RecognitionDateTo:    recognitionDateTo,
		FieldPurchaser:       fieldPurchaser,
		PurchasePlanDateFrom: purchasePlanDateFrom,
		PurchasePlanDateTo:   purchasePlanDateTo,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(purchasePlans, total, page)

	return ctx.Serve(err)
}

func (h PurchasePlanHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var purchasePlan *dto.PurchasePlanResponse

	id := c.Param("id")

	purchasePlan, err = h.ServicePurchasePlan.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = purchasePlan

	return ctx.Serve(err)
}

func (h PurchasePlanHandler) GetSummary(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var session dto.UserResponse
	var siteID string
	if session, err = h.ServicesAuth.Session(ctx.Request().Context()); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	if session.Site.ID != "" {
		siteID = session.Site.ID
	}

	fieldPurchaser := ctx.GetParamString("field_purchaser")
	if session.MainRole == "Field Purchaser" {
		fieldPurchaser = session.EmployeeCode
	}

	var summaryPurchasePlan *dto.SummaryPurchasePlanResponse
	summaryPurchasePlan, err = h.ServicePurchasePlan.GetSummary(ctx.Request().Context(), &dto.PurchasePlanListRequest{
		SiteID:         siteID,
		FieldPurchaser: fieldPurchaser,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = summaryPurchasePlan

	return ctx.Serve(err)
}

func (h PurchasePlanHandler) Assign(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchasePlanRequestAssign
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	id := c.Param("id")

	req.Session, _ = service.NewAuthService().Session(ctx.Request().Context())

	var res *dto.PurchasePlanResponse
	res, err = h.ServicePurchasePlan.Assign(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}

func (h PurchasePlanHandler) CancelAssign(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := c.Param("id")

	var res *dto.PurchasePlanResponse
	res, err = h.ServicePurchasePlan.CancelAssign(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}
