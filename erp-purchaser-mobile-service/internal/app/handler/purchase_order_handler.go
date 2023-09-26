package handler

import (
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type PurchaseOrderHandler struct {
	Option               global.HandlerOptions
	ServicePurchaseOrder service.IPurchaseOrderService
	ServicesAuth         service.IAuthService
}

// URLMapping declare endpoint with handler function.
func (h *PurchaseOrderHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicePurchaseOrder = service.NewPurchaseOrderService()

	h.ServicesAuth = service.NewAuthService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized("purchaser_app"))
	r.GET("/:id", h.Detail, cMiddleware.Authorized("purchaser_app"))
	r.POST("", h.Create, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/:id", h.Update, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/assign/:id", h.Assign, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/print/:id", h.Print, cMiddleware.Authorized("purchaser_app"))
	r.POST("/signature/:id", h.Signature, cMiddleware.Authorized("purchaser_app"))
	r.PUT("/cancel/:id", h.Cancel, cMiddleware.Authorized("purchaser_app"))
}

func (h PurchaseOrderHandler) Get(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var session dto.UserResponse
	var employeeCode string
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
	recognitionDateFrom := ctx.GetParamDate("recognition_date_from")
	recognitionDateTo := ctx.GetParamDate("recognition_date_to")
	code := ctx.GetParamArrayString("code")
	isNotConsolidatedStr := ctx.GetParamString("is_not_consolidated")
	purchasePlanIdStr := ctx.GetParamString("purchase_plan_id")
	prpcsno := ctx.GetParamString("prp_cs_no")
	employeeCode = ctx.GetParamString("field_purchaser")

	if session.MainRole == "Field Purchaser" {
		employeeCode = session.EmployeeCode
	}

	var isNotConsolidated bool
	if isNotConsolidatedStr != "" {
		isNotConsolidated, err = strconv.ParseBool(isNotConsolidatedStr)
		if err != nil {
			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
			return
		}
	}

	var purchaseOrders []*dto.PurchaseOrderResponse
	var total int64
	purchaseOrders, total, err = h.ServicePurchaseOrder.Get(ctx.Request().Context(), &dto.PurchaseOrderListRequest{
		Limit:               int32(page.Limit),
		Offset:              int32(page.Offset),
		Status:              int32(status),
		Search:              search,
		OrderBy:             orderBy,
		RecognitionDateFrom: timex.ToStartTime(recognitionDateFrom),
		RecognitionDateTo:   timex.ToLastTime(recognitionDateTo),
		Code:                code,
		IsNotConsolidated:   isNotConsolidated,
		PurchasePlanID:      purchasePlanIdStr,
		EmployeeCode:        employeeCode,
		Site:                session.Site.ID,
		PrpCsNo:             prpcsno,
	})
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(purchaseOrders, total, page)

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var purchaseOrder *dto.PurchaseOrderResponse

	id := c.Param("id")

	purchaseOrder, err = h.ServicePurchaseOrder.GetByID(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = purchaseOrder

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Create(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchaseOrderRequestCreate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}
	fmt.Print("payload ")
	fmt.Println(">>>>>> ", req)
	var res *dto.PurchaseOrderResponse
	res, err = h.ServicePurchaseOrder.Create(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Update(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchaseOrderRequestUpdate
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	id := c.Param("id")

	var res *dto.PurchaseOrderResponse
	res, err = h.ServicePurchaseOrder.Update(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Assign(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchaseOrderRequestAssign
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	id := c.Param("id")

	var res *dto.PurchaseOrderResponse
	res, err = h.ServicePurchaseOrder.Assign(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Print(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	id := c.Param("id")

	var res *dto.PurchaseOrderResponse
	res, err = h.ServicePurchaseOrder.Print(ctx.Request().Context(), id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Signature(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchaseOrderRequestSignature
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	id := c.Param("id")

	err = h.ServicePurchaseOrder.Signature(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h PurchaseOrderHandler) Cancel(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var req dto.PurchaseOrderRequestCancel
	if err = ctx.Bind(&req); err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	id := c.Param("id")

	var res *dto.PurchaseOrderResponse
	res, err = h.ServicePurchaseOrder.Cancel(ctx.Request().Context(), req, id)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	ctx.ResponseData = res

	return ctx.Serve(err)
}
