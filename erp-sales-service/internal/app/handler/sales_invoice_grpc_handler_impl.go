package handler

import (
	context "context"

	// bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	salesService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SalesGrpcHandler) GetSalesInvoiceGPMobileList(ctx context.Context, req *salesService.GetSalesInvoiceGPMobileListRequest) (res *salesService.GetSalesInvoiceGPMobileListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetSalesOrderList")
	defer span.End()
	//var total int64
	res1, err := h.ServiceSalesInvoice.GetSalesInvoiceGPMobileList(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	res = &res1
	return
}
