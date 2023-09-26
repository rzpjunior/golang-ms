package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetTransactionDetailGPList(ctx context.Context, req *bridgeService.GetTransactionDetailGPListRequest) (res *bridgeService.GetTransactionDetailGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransactionDetailGPList")
	defer span.End()

	res, err = h.ServicesTransactionDetail.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetTransactionDetailGPDetail(ctx context.Context, req *bridgeService.GetTransactionDetailGPDetailRequest) (res *bridgeService.GetTransactionDetailGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetTransactionDetailGPDetail")
	defer span.End()

	res, err = h.ServicesTransactionDetail.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
