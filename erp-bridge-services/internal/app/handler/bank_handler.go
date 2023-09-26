package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/labstack/echo/v4"
)

type BankHandler struct {
	Option       global.HandlerOptions
	ServicesBank service.IBankService
}

// URLMapping implements router.RouteHandlers
func (h *BankHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesBank = service.NewBankService()

	cMiddleware := middleware.NewMiddleware()
	r.GET("", h.Get, cMiddleware.Authorized())
	r.GET("/:id", h.Detail, cMiddleware.Authorized())
}

func (h BankHandler) Get(c echo.Context) (err error) {
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
	status := ctx.GetParamArrayInt("status")

	statusArr := make([]int32, len(status))

	for i, value := range status {
		statusArr[i] = int32(value)
	}

	var Bankes []dto.BankResponse
	var req *pb.GetBankListRequest
	req = &pb.GetBankListRequest{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Start),
		Status:  statusArr,
		Search:  search,
		OrderBy: orderBy,
	}

	var total int64
	Bankes, total, err = h.ServicesBank.Get(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.DataList(Bankes, total, page)

	return ctx.Serve(err)
}

func (h BankHandler) Detail(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	var Bank dto.BankResponse

	var id int64
	if id, err = ctx.GetParamID(); err != nil {
		err = edenlabs.ErrorValidation("id", "Invalid id")
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	var req *pb.GetBankDetailRequest
	req = &pb.GetBankDetailRequest{
		Id: id,
	}

	Bank, err = h.ServicesBank.GetDetail(ctx.Request().Context(), req)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return ctx.Serve(err)
	}

	ctx.ResponseData = Bank

	return ctx.Serve(err)
}
