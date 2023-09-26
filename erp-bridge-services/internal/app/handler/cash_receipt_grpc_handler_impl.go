package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetCashReceiptList(ctx context.Context, req *bridgeService.GetCashReceiptListRequest) (res *bridgeService.GetCashReceiptListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPList")
	defer span.End()

	res, err = h.ServicesCashReceipt.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateCashReceipt(ctx context.Context, req *bridgeService.CreateCashReceiptRequest) (res *bridgeService.CreateCashReceiptResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressGPList")
	defer span.End()

	res, err = h.ServicesCashReceipt.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
